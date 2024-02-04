package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"health": "hello world",
	})
}
