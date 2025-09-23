package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /api/v1/traq/users/:id
// traQ APIのユーザー詳細をプロキシ
func (h *Handler) GetTraqUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	token := getTokenFromCookie(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	url := "https://q.trap.jp/api/v3/users/" + id
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, url, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build request").SetInternal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "failed to request traQ").SetInternal(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warn: failed to close response body: %v", cerr)
		}
	}()

	return h.proxyResponse(c, resp)
}
