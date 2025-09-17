package handler

import (
	"net/http"

	"github.com/traP-jp/1m25_10/backend/internal/handler/middleware"

	"github.com/traP-jp/1m25_10/backend/internal/repository"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo   repository.Repository
	client *http.Client
}

// New creates a Handler. If client is nil, http.DefaultClient will be used.
func New(repo repository.Repository, client *http.Client) *Handler {
	if client == nil {
		client = http.DefaultClient
	}

	return &Handler{
		repo:   repo,
		client: client,
	}
}

func (h *Handler) SetupAppRoutes(api *echo.Group) {
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
		albumAPI.DELETE("/:id", h.DeleteAlbum, middleware.UsernameProvider)
		// Prefer PATCH for partial updates; keep PUT for backward compatibility
		albumAPI.PATCH("/:id", h.UpdateAlbum, middleware.UsernameProvider)
		albumAPI.PUT("/:id", h.UpdateAlbum, middleware.UsernameProvider)
	}

	// images API
	imagesAPI := api.Group("/images")
	{
		imagesAPI.GET("", h.GetTraqMessagesSearchImages)
		imagesAPI.GET("/:id", h.GetLatestMessageByImageID)
	}

}

// SetupAuthRoutes は `/api/auth` にマウントされる Auth 専用ルートを登録します。
// 引数の `authGroup` は既に `/api/auth` のグループであることを想定します。
func (h *Handler) SetupAuthRoutes(authGroup *echo.Group) {
	authGroup.GET("/request", h.AuthRequest)
	authGroup.GET("/callback", h.AuthCallback)
	authGroup.GET("/me", h.AuthMe)
	authGroup.POST("/logout", h.AuthLogout)
}

// SetupTraqRoutes は `/api/v1/traq` にマウントされる traQ プロキシ専用ルートを登録します。
// 引数の `traqGroup` は既に `/api/v1/traq` のグループであることを想定します。
func (h *Handler) SetupTraqRoutes(traqGroup *echo.Group) {
	filesGroup := traqGroup.Group("/files")
	{
		filesGroup.GET("/:uuid", h.GetTraqFile)
		filesGroup.GET("/:uuid/thumbnail", h.GetTraqFileThumbnail)
	}
	// messages search endpoints
	traqGroup.GET("/messages", h.GetTraqMessagesSearch)
}
