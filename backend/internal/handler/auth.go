package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_10/backend/pkg/config"
)

const (
	cookieTokenKey    = "traq-auth-token"
	cookieStateKey    = "traq-auth-state"
	cookieVerifierKey = "traq-auth-code-verifier"
	cookieCallbackKey = "traq-auth-callback"
	authBase          = "https://q.trap.jp/api/v3/oauth2"
)

// GET /api/auth/request
// query: callback (optional, relative path "/...")
func (h *Handler) AuthRequest(c echo.Context) error {
	clientID := config.TraqOAuthClientID()
	redirectURI := config.TraqOAuthRedirectURI()
	if clientID == "" || redirectURI == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "OAuth not configured")
	}

	// state と code_verifier を生成
	state, err := randString(32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate state").SetInternal(err)
	}
	codeVerifier, err := randString(64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate code_verifier").SetInternal(err)
	}
	codeChallenge := s256(codeVerifier)

	// 一時Cookieへ保存（短寿命）
	setTempCookie(c, cookieStateKey, state, 10*time.Minute)
	setTempCookie(c, cookieVerifierKey, codeVerifier, 10*time.Minute)

	cb := c.QueryParam("callback")
	if cb != "" && strings.HasPrefix(cb, "/") {
		setTempCookie(c, cookieCallbackKey, cb, 15*time.Minute)
	}

	// 認可エンドポイントへリダイレクト
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", clientID)
	v.Set("redirect_uri", redirectURI)
	v.Set("state", state)
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", "S256")
	// v.Set("scope", "read") // 必要に応じて

	authURL := authBase + "/authorize?" + v.Encode()
	return c.Redirect(http.StatusFound, authURL)
}

// GET /api/auth/callback?code=...&state=...
func (h *Handler) AuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	rstate := c.QueryParam("state")
	if code == "" || rstate == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing code/state")
	}

	stateCookie, err := c.Cookie(cookieStateKey)
	if err != nil || stateCookie.Value == "" || stateCookie.Value != rstate {
		return echo.NewHTTPError(http.StatusBadRequest, "state mismatch")
	}
	verifierCookie, err := c.Cookie(cookieVerifierKey)
	if err != nil || verifierCookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "verifier missing")
	}

	// トークン交換
	token, expiresIn, err := exchangeToken(code, verifierCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "token exchange failed").SetInternal(err)
	}

	// アクセストークンをCookieに保存
	setAuthCookie(c, cookieTokenKey, token, time.Duration(expiresIn)*time.Second)

	// 一時Cookieを削除
	delCookie(c, cookieStateKey)
	delCookie(c, cookieVerifierKey)

	// callbackへリダイレクト
	cb := "/"
	if cbCookie, err := c.Cookie(cookieCallbackKey); err == nil && cbCookie.Value != "" && strings.HasPrefix(cbCookie.Value, "/") {
		cb = cbCookie.Value
	}
	delCookie(c, cookieCallbackKey)
	// 開発時はフロントの別ポートへ戻す
	if fb := config.FrontendBaseURL(); fb != "" {
		base := strings.TrimRight(fb, "/")
		return c.Redirect(http.StatusFound, base+cb)
	}
	return c.Redirect(http.StatusFound, cb)
}

// GET /api/auth/me
func (h *Handler) AuthMe(c echo.Context) error {
	token := getTokenFromCookie(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	// traQの /users/me を呼んでusernameを返す
	u, err := fetchTraqMe(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}
	return c.JSON(http.StatusOK, u)
}

// POST /api/auth/logout
func (h *Handler) AuthLogout(c echo.Context) error {
	delCookie(c, cookieTokenKey)
	return c.NoContent(http.StatusNoContent)
}

// ========== helpers ==========
func randString(n int) (string, error) {
	// Base64URL (no padding) で、均一分布かつURLセーフな文字列を生成する。
	// 望みの長さ n に対して、原始バイト長は ceil(n * 3 / 4)。不足時は切り詰める。
	if n <= 0 {
		return "", nil
	}
	bytesLen := (n*3 + 3) / 4
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) >= n {
		return s[:n], nil
	}
	// 通常は発生しないが、不足した場合は追加生成して充足させる
	extra := make([]byte, (n-len(s))*3/4+1)
	if _, err := rand.Read(extra); err != nil {
		return "", err
	}
	s += base64.RawURLEncoding.EncodeToString(extra)
	if len(s) < n {
		return "", errors.New("failed to generate enough random data")
	}
	return s[:n], nil
}

func s256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func setTempCookie(c echo.Context, name, val string, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(maxAge / time.Second),
	}
	c.SetCookie(cookie)
}

func setAuthCookie(c echo.Context, name, val string, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(maxAge / time.Second),
	}
	c.SetCookie(cookie)
}

func delCookie(c echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(cookie)
}

func getTokenFromCookie(c echo.Context) string {
	if ck, err := c.Cookie(cookieTokenKey); err == nil {
		return ck.Value
	}
	return ""
}

// ------------- HTTP calls to traQ -------------
type meResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func exchangeToken(code, codeVerifier string) (accessToken string, expiresIn int64, err error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", config.TraqOAuthClientID())
	if cs := config.TraqOAuthClientSecret(); cs != "" {
		form.Set("client_secret", cs)
	}
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	if ru := config.TraqOAuthRedirectURI(); ru != "" {
		form.Set("redirect_uri", ru)
	}

	req, err := http.NewRequest(http.MethodPost, authBase+"/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", 0, err
	}
	// x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", 0, errors.New("token endpoint error (" + resp.Status + "): " + string(b))
	}
	// 最小限のparse（access_tokenとexpires_inのみ）
	type tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	var tr tokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", 0, err
	}
	return tr.AccessToken, tr.ExpiresIn, nil
}

func fetchTraqMe(token string) (*meResponse, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://q.trap.jp/api/v3/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New("me error: " + string(b))
	}
	var me meResponse
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return nil, err
	}
	return &me, nil
}
