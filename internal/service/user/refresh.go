package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Refresh(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
