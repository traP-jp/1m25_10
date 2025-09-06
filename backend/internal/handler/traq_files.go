package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /api/v1/traq/files/{uuid}
// traQ APIのファイル本体を取得してプロキシする
func (h *Handler) GetTraqFile(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "uuid is required")
	}

	// CookieからtraQ認証トークンを取得
	token := getTokenFromCookie(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	// traQ APIにリクエストを送信
	resp, err := h.proxyTraqFileRequest("https://q.trap.jp/api/v3/files/"+uuid, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch file from traQ").SetInternal(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warn: failed to close response body: %v", cerr)
		}
	}()

	// traQ APIのレスポンスをそのままクライアントに返す
	return h.proxyResponse(c, resp)
}

// GET /api/v1/traq/files/{uuid}/thumbnail
// traQ APIのファイルサムネイルを取得してプロキシする
func (h *Handler) GetTraqFileThumbnail(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "uuid is required")
	}

	// CookieからtraQ認証トークンを取得
	token := getTokenFromCookie(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	// traQ APIにリクエストを送信
	resp, err := h.proxyTraqFileRequest("https://q.trap.jp/api/v3/files/"+uuid+"/thumbnail", token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch thumbnail from traQ").SetInternal(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warn: failed to close response body: %v", cerr)
		}
	}()

	// traQ APIのレスポンスをそのままクライアントに返す
	return h.proxyResponse(c, resp)
}

// traQ APIへのHTTPリクエストを送信する共通関数
func (h *Handler) proxyTraqFileRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// AuthorizationヘッダーにtraQトークンを設定
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// traQ APIのレスポンスをクライアントにそのまま転送する共通関数
func (h *Handler) proxyResponse(c echo.Context, resp *http.Response) error {
	// ステータスコードを設定
	c.Response().Status = resp.StatusCode

	// 重要なヘッダーをコピー
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		c.Response().Header().Set("Content-Type", contentType)
	}
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Response().Header().Set("Content-Length", contentLength)
	}
	if cacheControl := resp.Header.Get("Cache-Control"); cacheControl != "" {
		c.Response().Header().Set("Cache-Control", cacheControl)
	}
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.Response().Header().Set("ETag", etag)
	}
	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		c.Response().Header().Set("Last-Modified", lastModified)
	}

	// レスポンスボディをそのままコピー
	// エラーレスポンスも含めて、すべてのレスポンスをそのまま転送
	_, err := io.Copy(c.Response().Writer, resp.Body)
	return err
}
