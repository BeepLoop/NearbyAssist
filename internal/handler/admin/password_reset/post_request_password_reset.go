package passwordreset_handler

import (
	"nearbyassist/internal/service/activitylog"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *passwordResetHandler) RequestPasswordReset(c echo.Context) error {
	username := c.FormValue("username")

	admin, _ := h.passwordResetService.GetAdminFromUsername(username)

	if err := h.passwordResetService.RequestPasswordReset(username); err != nil {
		if strings.Contains(err.Error(), passwordreset_service.ERR_INVALID_USERNAME) {
			if err := utils.SetFlashMessage(c, "error", "Invalid username"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/?error=invalid_username")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Error requesting password reset"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/?error=request_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	if err := utils.SetFlashMessage(c, "success", "Request submitted"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/?success=request_submitted")
	}

	activity := activitylog.Input{
		AdminId: admin.Id,
		Action:  activitylog.ACTION_PASSWORD_RESET_REQUEST,
	}
	activitylog.MustGetInstance().Create(activity)
	sse.New().IncreasePasswordResetRequest()

	return c.Redirect(http.StatusSeeOther, "/")
}
