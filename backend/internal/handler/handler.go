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
}
