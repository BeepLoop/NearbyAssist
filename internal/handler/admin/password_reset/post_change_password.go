package passwordreset_handler

import (
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *passwordResetHandler) PostChangePassword(c echo.Context) error {
	username := c.FormValue("username")
	oldPassword := c.FormValue("oldPassword")
	password := c.FormValue("password")
	confirmationPassword := c.FormValue("confirmationPassword")

	if err := h.passwordResetService.ChangePassword(username, oldPassword, password, confirmationPassword); err != nil {
		if strings.Contains(err.Error(), passwordreset_service.ERR_MISMATCHING_PASSWORD) {
			if err := utils.SetFlashMessage(c, "error", "Confirmation password mismatch"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/reset/cp?error=mismatch_password_error")
			}
		} else if strings.Contains(err.Error(), passwordreset_service.ERR_INVALID_CREDENTIALS) {
			if err := utils.SetFlashMessage(c, "error", "Invalid credentials"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/reset/cp?error=invalid_credentials")
			}
		} else if strings.Contains(err.Error(), passwordreset_service.ERR_WEAK_PASSWORD) {
			if err := utils.SetFlashMessage(c, "error", "Provided password is too weak"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/reset/cp?error=weak_password")
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
