package auth

import (
	"context"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/auth"

	"github.com/labstack/echo/v4"
)

func (h *authHandler) GetLogin(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	page := pages.Login(flash)
	return page.Render(context.Background(), c.Response().Writer)
}
