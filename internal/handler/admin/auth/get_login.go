package auth

import (
	"context"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *authHandler) GetLogin(c echo.Context) error {
	sess, err := utils.GetAdminFromSession(c)
	if err == nil && sess != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	flash, _, err := utils.RetrieveFlashMessage(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/?error=session_error")
	}

	page := pages.Login(flash)
	return page.Render(context.Background(), c.Response().Writer)
}
