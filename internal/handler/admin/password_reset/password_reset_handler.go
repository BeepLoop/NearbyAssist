package passwordreset_handler

import (
	"context"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type passwordResetHandler struct {
	passwordResetService *passwordreset_service.Service
}

func NewHandler(passwordResetService *passwordreset_service.Service) *passwordResetHandler {
	return &passwordResetHandler{
		passwordResetService: passwordResetService,
	}
}

func (h *passwordResetHandler) RequestPasswordReset(c echo.Context) error {
	username := c.FormValue("username")

	if err := h.passwordResetService.RequestPasswordReset(username); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Request failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/?error=request_error")
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	if err := utils.SetFlashMessage(c, "success", "Request submitted"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/?success=request_submitted")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *passwordResetHandler) ChangePassword(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	page := pages.ChangePassword(flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *passwordResetHandler) PostChangePassword(c echo.Context) error {
	username := c.FormValue("username")
	oldPassword := c.FormValue("oldPassword")
	password := c.FormValue("password")
	confirmationPassword := c.FormValue("confirmationPassword")

	if err := h.passwordResetService.ChangePassword(username, oldPassword, password, confirmationPassword); err != nil {
		if strings.Contains(err.Error(), "password mismatch") {
			if err := utils.SetFlashMessage(c, "error", "Mismatching password"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/reset/cp?error=mismatch_password_error")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Password change failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/reset/cp?error=password_change_failed")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/reset/cp")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}
