package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", s.HealthHandler)

	return e
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"health": "hello world",
	})
}
