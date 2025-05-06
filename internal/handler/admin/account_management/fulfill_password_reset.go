package accountmanagement

import (
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) FufillResetRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	password := c.FormValue("password")
	confirmationUsername := c.FormValue("confirmationUsername")
	confirmationPassword := c.FormValue("confirmationPassword")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	input := passwordreset_service.ResetRequestInput{
		HandlerAdminId:       admin.Id,
		RequestId:            requestId,
		NewPassword:          password,
		ConfirmationUsername: confirmationUsername,
		ConfirmationPassword: confirmationPassword,
	}

	if err := h.passwordResetService.FulfillResetPassword(input); err != nil {
		if strings.Contains(err.Error(), passwordreset_service.ERR_SELF_RESETTING_PASSWORD) {
			if err := utils.SetFlashMessage(c, "error", "Resetting own password not allowed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=resetting_own_password_disallowed")
			}
		} else if strings.Contains(err.Error(), passwordreset_service.ERR_INVALID_CREDENTIALS) {
			if err := utils.SetFlashMessage(c, "error", "Invalid confirmation credentials"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=invalid_credentials")
			}
		} else if strings.Contains(err.Error(), passwordreset_service.ERR_WEAK_PASSWORD) {
			if err := utils.SetFlashMessage(c, "error", "Provided password is too weak"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=weak_password")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Password reset failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=reset_failed")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	if err := utils.SetFlashMessage(c, "success", "Request success"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?success=password_change_success")
	}

	sse.New().PasswordResetRequest--

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
}
