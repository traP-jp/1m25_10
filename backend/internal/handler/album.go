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
	creatorIdStr := c.QueryParam("creator_id")
	var creatorId *uuid.UUID
	if creatorIdStr != "" {
		creatorIdParsed, err := uuid.Parse(creatorIdStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid creator ID")
		}
		creatorId = &creatorIdParsed
	}
	beforeDateStr := c.QueryParam("before_date")
	var beforeDate *time.Time
	if beforeDateStr != "" {
		beforeDateParsed, err := time.Parse(time.RFC3339, beforeDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid before date")
		}
		beforeDate = &beforeDateParsed
	}
	afterDateStr := c.QueryParam("after_date")
	var afterDate *time.Time
	if afterDateStr != "" {
		afterDateParsed, err := time.Parse(time.RFC3339, afterDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid after date")
		}
		afterDate = &afterDateParsed
	}
	limitStr := c.QueryParam("limit")
	var limit *int
	if limitStr != "" {
		limitParsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid limit")
		}
		limit = &limitParsed
	}
	offsetStr := c.QueryParam("offset")
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
