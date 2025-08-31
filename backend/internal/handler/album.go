package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)


func (h *Handler) GetAlbum(c echo.Context) error {
	albumID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid album ID")
	}
	album, err := h.repo.GetAlbum(c.Request().Context(), albumID)
	if err != nil {
		if err.Error() ==  "album not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Album not found")
		}
	}
	return c.JSON(http.StatusOK, album)
}
