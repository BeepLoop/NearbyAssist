package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HandleVersionOneRoutes(r *echo.Group) {

	r.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

}
