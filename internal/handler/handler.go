package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func (h *Handler) HealthCheck(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
