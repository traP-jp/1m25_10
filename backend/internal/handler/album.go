package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/traP-jp/1m25_10/backend/internal/domain"
	"github.com/traP-jp/1m25_10/backend/internal/handler/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GET /api/v1/albums
func (h *Handler) GetAlbums(c echo.Context) error {
	creatorIdStr := c.QueryParam("creator")
	var creatorId *string
	if creatorIdStr != "" {
		creatorId = &creatorIdStr
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

	return c.JSON(http.StatusOK, albums)
}

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

// POST /api/v1/albums
func (h *Handler) PostAlbum(c echo.Context) error {
	req := new(struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	})
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	creator, ok := c.Get(middleware.UsernameKey).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	images := make([]uuid.UUID, 0, len(req.Images))
	for _, s := range req.Images {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid image id: %s", s)).SetInternal(err)
		}
		images = append(images, id)
	}

	params := domain.PostAlbumParams{
		Title:       req.Title,
		Description: req.Description,
		Creator:     creator,
		Images:      images,
	}

	album, err := h.repo.PostAlbum(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create album")
	}
	return c.JSON(http.StatusCreated, album)

}
