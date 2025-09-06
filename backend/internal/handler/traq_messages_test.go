package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type rtFunc func(*http.Request) *http.Response

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r), nil }

// Test that searchTraqMessages builds query and headers correctly and returns body on 200
func TestSearchTraqMessages_OK(t *testing.T) {
	// Prepare echo context with cookie
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
	// set auth cookie used by getTokenFromCookie
	cookie := &http.Cookie{Name: cookieTokenKey, Value: "test-token"}
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock HTTP client
	var capturedURL string
	var capturedAuth string
	client := &http.Client{Transport: rtFunc(func(r *http.Request) *http.Response {
		capturedURL = r.URL.String()
		capturedAuth = r.Header.Get("Authorization")
		body := io.NopCloser(strings.NewReader(`{"ok":true}`))
		return &http.Response{StatusCode: 200, Body: body, Header: make(http.Header)}
	})}

	h := &Handler{client: client}

	hasImage := true
	limit := 5
	p := &traqMessageSearchParams{
		Word:     "cat",
		To:       []string{"u1", "u2"},
		From:     []string{"f1"},
		HasImage: &hasImage,
		Limit:    &limit,
		Sort:     "createdAt",
	}

	body, status, err := h.searchTraqMessages(c, p)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if status != 200 {
		t.Fatalf("unexpected status: %d", status)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", string(body))
	}

	if !strings.Contains(capturedURL, "/api/v3/messages?") {
		t.Fatalf("unexpected URL: %s", capturedURL)
	}
	// Query checks
	for _, want := range []string{"word=cat", "to=u1", "to=u2", "from=f1", "hasImage=true", "limit=5", "sort=createdAt"} {
		if !strings.Contains(capturedURL, want) {
			t.Errorf("missing query param %q in %s", want, capturedURL)
		}
	}
	// Header check
	if capturedAuth != "Bearer test-token" {
		t.Fatalf("unexpected Authorization header: %s", capturedAuth)
	}
}

func TestExtractUUIDsFromContent(t *testing.T) {
	content := `text https://q.trap.jp/files/01990957-7a91-79c8-bb64-38766a3411e6 and also https://q.trap.jp/files/550e8400-e29b-41d4-a716-446655440000`
	got := extractUUIDsFromContent(content)
	if len(got) != 2 {
		t.Fatalf("expected 2 uuids, got %d: %v", len(got), got)
	}
}

func TestSearchTraqImagesUUIDs_ParseOnly(t *testing.T) {
	// Build handler with client returning canned JSON
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
	req.AddCookie(&http.Cookie{Name: cookieTokenKey, Value: "t"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	sample := `{"totalHits":2,"hits":[{"content":"see https://q.trap.jp/files/550e8400-e29b-41d4-a716-446655440000"},{"content":"two: https://q.trap.jp/files/01990957-7a91-79c8-bb64-38766a3411e6 https://q.trap.jp/files/550e8400-e29b-41d4-a716-446655440001"}]}`
	client := &http.Client{Transport: rtFunc(func(r *http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(sample)), Header: make(http.Header)}
	})}
	h := &Handler{client: client}

	total, uuids, err := h.searchTraqImagesUUIDs(c, &traqMessageSearchParams{Word: "テスト"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if total != 2 {
		t.Fatalf("unexpected totalHits: %d", total)
	}
	if len(uuids) != 3 {
		t.Fatalf("expected 3 uuids, got %d: %v", len(uuids), uuids)
	}
}
