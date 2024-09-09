package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "me")
}
