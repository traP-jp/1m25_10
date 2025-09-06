package handler

import (
	"github.com/traP-jp/1m25_10/backend/internal/handler/middleware"
	"github.com/traP-jp/1m25_10/backend/internal/repository"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) SetupRoutes(api *echo.Group) {
	// ping API
	pingAPI := api.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// album API
	albumAPI := api.Group("/albums")
	{
		albumAPI.GET("", h.GetAlbums)
		albumAPI.GET("/:id", h.GetAlbum)
		albumAPI.POST("", h.PostAlbum, middleware.UsernameProvider)
	}
}
