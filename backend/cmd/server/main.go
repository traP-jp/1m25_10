package main

import (
	"github.com/traP-jp/1m25_10/backend/cmd/server/server"
	"github.com/traP-jp/1m25_10/backend/pkg/config"
	"github.com/traP-jp/1m25_10/backend/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to and migrate database
	db, err := database.Setup(config.MySQL())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			e.Logger.Errorf("failed to close DB: %v", err)
		}
	}()

	s := server.Inject(db)

	v1API := e.Group("/api/v1")
	s.SetupRoutes(v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
