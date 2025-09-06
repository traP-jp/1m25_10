package server

import (
	"net/http"
	"time"

	"github.com/traP-jp/1m25_10/backend/internal/handler"
	"github.com/traP-jp/1m25_10/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	handler *handler.Handler
}

func Inject(db *sqlx.DB) *Server {
	repo := repository.New(db)

	// Create an HTTP client with a reasonable timeout for external calls
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	h := handler.New(repo, client)

	return &Server{
		handler: h,
	}
}

// ルートレベルのセットアップ
func (d *Server) SetupRoot(e *echo.Echo) {
	// top-level /api group
	api := e.Group("/api")

	// /api/auth
	authGroup := api.Group("/auth")
	d.handler.SetupAuthRoutes(authGroup)

	// /api/v1
	v1Group := api.Group("/v1")
	d.handler.SetupAppRoutes(v1Group)

	// /api/v1/traq
	traqGroup := v1Group.Group("/traq")
	d.handler.SetupTraqRoutes(traqGroup)
}
