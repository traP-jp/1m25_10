package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

// traqMessageSearchParams は traQ /api/v3/messages のクエリに対応する内部用パラメータです。
type traqMessageSearchParams struct {
	Word           string
	After          string // RFC3339想定 (serverはバリデーションしない)
	Before         string // RFC3339想定 (serverはバリデーションしない)
	In             string // channel uuid
	To             []string
	From           []string
	Citation       string // message uuid
	Bot            *bool
	HasURL         *bool
	HasAttachments *bool
	HasImage       *bool
	HasVideo       *bool
	HasAudio       *bool
	Limit          *int
	Offset         *int
	Sort           string // createdAt | -createdAt | updatedAt | -updatedAt
}

// searchTraqMessages は traQ のメッセージ検索APIをそのまま実行し、レスポンスボディ(JSON)を返します。
// HTTP経由でそのまま返すのではなく、後続処理用の生データ取得に利用します。
// 返り値: body, statusCode, error
func (h *Handler) searchTraqMessages(c echo.Context, p *traqMessageSearchParams) ([]byte, int, error) {
	token := getTokenFromCookie(c)
	if token == "" {
		return nil, http.StatusUnauthorized, echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	if p == nil {
		p = &traqMessageSearchParams{}
	}

	endpoint, _ := url.Parse("https://q.trap.jp/api/v3/messages")
	q := endpoint.Query()

	if p.Word != "" {
		q.Set("word", p.Word)
	}
	if p.After != "" {
		q.Set("after", p.After)
	}
	if p.Before != "" {
		q.Set("before", p.Before)
	}
	if p.In != "" {
		q.Set("in", p.In)
	}
	for _, v := range p.To {
		if v != "" {
			q.Add("to", v)
		}
	}
	for _, v := range p.From {
		if v != "" {
			q.Add("from", v)
		}
	}
	if p.Citation != "" {
		q.Set("citation", p.Citation)
	}
	if p.Bot != nil {
		q.Set("bot", strconv.FormatBool(*p.Bot))
	}
	if p.HasURL != nil {
		q.Set("hasURL", strconv.FormatBool(*p.HasURL))
	}
	if p.HasAttachments != nil {
		q.Set("hasAttachments", strconv.FormatBool(*p.HasAttachments))
	}
	if p.HasImage != nil {
		q.Set("hasImage", strconv.FormatBool(*p.HasImage))
	}
	if p.HasVideo != nil {
		q.Set("hasVideo", strconv.FormatBool(*p.HasVideo))
	}
	if p.HasAudio != nil {
		q.Set("hasAudio", strconv.FormatBool(*p.HasAudio))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if p.Sort != "" {
		q.Set("sort", p.Sort)
	}

	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warn: failed to close response body: %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 生のエラーボディを付与
		return body, resp.StatusCode, fmt.Errorf("traQ messages search failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}

// ---------- 画像UUID抽出用の内部処理 ----------

// traqMessagesRawResponse は最小限のフィールドのみ定義
type traqMessagesRawResponse struct {
	TotalHits int `json:"totalHits"`
	Hits      []struct {
		Content string `json:"content"`
	} `json:"hits"`
}

var fileUUIDRe = regexp.MustCompile(`https?://q\.trap\.jp/files/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

// メッセージ本文からファイルUUIDを全件抽出する
func extractUUIDsFromContent(content string) []string {
	matches := fileUUIDRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}
	res := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) >= 2 {
			res = append(res, m[1])
		}
	}
	return res
}

// hasImage=true 固定で traQ 検索を行い、totalHits と content から抽出した画像UUIDの配列を返す。
func (h *Handler) searchTraqImagesUUIDs(c echo.Context, p *traqMessageSearchParams) (int, []string, error) {
	// callerからのHasImage指定は無視してtrue固定
	has := true
	if p == nil {
		p = &traqMessageSearchParams{}
	}
	p.HasImage = &has

	body, status, err := h.searchTraqMessages(c, p)
	if err != nil {
		// 上位で扱いやすいように、traQのステータスも含める
		return 0, nil, fmt.Errorf("traQ search failed (status=%d): %w", status, err)
	}

	var raw traqMessagesRawResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return 0, nil, err
	}

	uuids := make([]string, 0)
	for _, h := range raw.Hits {
		if h.Content == "" {
			continue
		}
		found := extractUUIDsFromContent(h.Content)
		if len(found) > 0 {
			uuids = append(uuids, found...)
		}
	}
	return raw.TotalHits, uuids, nil
}

// traQ検索を行い、totalHits と抽出した画像UUID配列を返す。
func (h *Handler) GetTraqMessagesSearchImages(c echo.Context) error {
	params := &traqMessageSearchParams{
		Word:     c.QueryParam("word"),
		After:    c.QueryParam("after"),
		Before:   c.QueryParam("before"),
		In:       c.QueryParam("in"),
		Citation: c.QueryParam("citation"),
		Sort:     c.QueryParam("sort"),
	}
	if v := c.QueryParam("bot"); v != "" {
		b := v == "true" || v == "1"
		params.Bot = &b
	}
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.Limit = &n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.Offset = &n
		}
	}
	params.To = c.QueryParams()["to"]
	params.From = c.QueryParams()["from"]

	total, uuids, err := h.searchTraqImagesUUIDs(c, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// レスポンス整形
	out := map[string]interface{}{
		"totalHits": total,
		"hits":      uuids,
	}
	return c.JSON(http.StatusOK, out)
}

// 透過プロキシエンドポイント。
func (h *Handler) GetTraqMessagesSearch(c echo.Context) error {
	// クエリを構築
	params := &traqMessageSearchParams{
		Word:     c.QueryParam("word"),
		After:    c.QueryParam("after"),
		Before:   c.QueryParam("before"),
		In:       c.QueryParam("in"),
		Citation: c.QueryParam("citation"),
		Sort:     c.QueryParam("sort"),
	}
	// bool/intは存在チェックしてから
	if v := c.QueryParam("bot"); v != "" {
		b := v == "true" || v == "1"
		params.Bot = &b
	}
	if v := c.QueryParam("hasURL"); v != "" {
		b := v == "true" || v == "1"
		params.HasURL = &b
	}
	if v := c.QueryParam("hasAttachments"); v != "" {
		b := v == "true" || v == "1"
		params.HasAttachments = &b
	}
	if v := c.QueryParam("hasImage"); v != "" {
		b := v == "true" || v == "1"
		params.HasImage = &b
	}
	if v := c.QueryParam("hasVideo"); v != "" {
		b := v == "true" || v == "1"
		params.HasVideo = &b
	}
	if v := c.QueryParam("hasAudio"); v != "" {
		b := v == "true" || v == "1"
		params.HasAudio = &b
	}
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.Limit = &n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.Offset = &n
		}
	}
	// 多値クエリ
	params.To = c.QueryParams()["to"]
	params.From = c.QueryParams()["from"]

	// traQのエラーはそのまま返す（検証用途）
	body, status, _ := h.searchTraqMessages(c, params)
	return c.Blob(status, "application/json", body)
}

// GetLatestMessageByImageID
// GET /api/v1/images/:id
// 指定されたファイルUUIDを含む https://q.trap.jp/files/<uuid> を語句検索し、
// createdAt 昇順（古い順）に並べて最古の1件のヒットのみを返します。
func (h *Handler) GetLatestMessageByImageID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	word := "https://q.trap.jp/files/" + id
	limit := 1
	// 最古を取得するため昇順
	sort := "createdAt"
	params := &traqMessageSearchParams{
		Word:  word,
		Limit: &limit,
		Sort:  sort,
	}

	body, status, err := h.searchTraqMessages(c, params)
	if err != nil {
		// traQ 側のエラーはそのまま中継
		return c.Blob(status, "application/json", body)
	}

	// レスポンスから最古1件のhitのみ抽出（昇順ソートの先頭）
	var resp struct {
		Hits []json.RawMessage `json:"hits"`
	}
	if uerr := json.Unmarshal(body, &resp); uerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse traQ response").SetInternal(uerr)
	}
	if len(resp.Hits) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no message found for the given image id")
	}
	// 先頭(最古)のみ返す
	return c.Blob(http.StatusOK, "application/json", resp.Hits[0])
}
