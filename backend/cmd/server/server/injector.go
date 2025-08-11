package server

import (
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
	h := handler.New(repo)

	return &Server{
		handler: h,
	}
}

func (d *Server) SetupRoutes(g *echo.Group) {
	// TODO: handler.SetupRoutesを呼び出す or 直接書く？
	d.handler.SetupRoutes(g)
}
