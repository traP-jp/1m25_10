package handler

import (
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

func (h *Handler) SetupAppRoutes(api *echo.Group) {
	// ping API
	pingAPI := api.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user API
	userAPI := api.Group("/users")
	{
		userAPI.GET("", h.GetUsers)
		userAPI.POST("", h.CreateUser)
		userAPI.GET("/:userID", h.GetUser)
	}

	// album API
	albumAPI := api.Group("/albums")
	{
		albumAPI.GET("", h.GetAlbums)
		albumAPI.GET("/:id", h.GetAlbum)
	}
}

// SetupAuthRoutes は /api 直下にマウントされるAuth専用ルートを登録します。
func (h *Handler) SetupAuthRoutes(api *echo.Group) {
	authAPI := api.Group("/auth")
	{
		authAPI.GET("/request", h.AuthRequest)
		authAPI.GET("/callback", h.AuthCallback)
		authAPI.GET("/me", h.AuthMe)
		authAPI.POST("/logout", h.AuthLogout)
	}
}
