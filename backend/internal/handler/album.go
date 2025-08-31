package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/traP-jp/1m25_10/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// スキーマ定義
type (
	GetAlbumsResponse []domain.AlbumItem
)

// GET /api/v1/albums
func (h *Handler) GetAlbums(c echo.Context) error {
	creatorIdStr := c.Param("creator_id")
	var creatorId *uuid.UUID
	if creatorIdStr != "" {
		creatorIdParsed, err := uuid.Parse(creatorIdStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid creator ID")
		}
		creatorId = &creatorIdParsed
	}
	beforeDateStr := c.Param("before_date")
	var beforeDate *time.Time
	if beforeDateStr != "" {
		layout := "0000-01-01T00:00:00.000000Z"
		beforeDateParsed, err := time.Parse(layout, beforeDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid before date")
		}
		beforeDate = &beforeDateParsed
	}
	afterDateStr := c.Param("after_date")
	var afterDate *time.Time
	if afterDateStr != "" {
		layout := "0000-01-01T00:00:00.000000Z"
		afterDateParsed, err := time.Parse(layout, afterDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid after date")
		}
		afterDate = &afterDateParsed
	}
	limitStr := c.Param("limit")
	var limit *int
	if limitStr != "" {
		limitParsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid limit")
		}
		limit = &limitParsed
	}
	offsetStr := c.Param("offset")
	var offset *int
	if offsetStr != "" {
		offsetParsed, err := strconv.Atoi(offsetStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid offset")
		}
		offset = &offsetParsed
	}
	albums, err := h.repo.GetAlbums(c.Request().Context(), domain.AlbumFilter{
		CreatorID:  creatorId,
		BeforeDate: beforeDate,
		AfterDate:  afterDate,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := make(GetAlbumsResponse, len(albums))
	for i, album := range albums {
		response[i] = domain.AlbumItem{
			Id:      album.Id,
			Title:   album.Title,
			Creator: album.Creator,
		}
	}

	return c.JSON(http.StatusOK, response)
}
