package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_10/backend/internal/domain"
)


func (h *Handler) GetAlbum(c echo.Context) error {
	albumID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid album ID")
	}
	album, err := h.repo.GetAlbum(c.Request().Context(), albumID)
	if err != nil {
		if err == domain.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Album not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve album")
	}
	return c.JSON(http.StatusOK, album)
}
