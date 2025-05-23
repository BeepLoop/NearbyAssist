package auth

import (
	"context"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *authHandler) GetLogin(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	admin, _ := utils.GetAdminFromSession(c)
	if admin != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	page := pages.Login(flash)
	return page.Render(context.Background(), c.Response().Writer)
}
