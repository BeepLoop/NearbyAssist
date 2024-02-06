package server

import (
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HandleRegister(c echo.Context) error {
	u := new(types.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request data",
		})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request data",
		})
	}

	return c.JSON(http.StatusOK, u)
}
