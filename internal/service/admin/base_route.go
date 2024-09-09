package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AdminService) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "admin")
}
