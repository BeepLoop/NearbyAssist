package accountmanagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) RejectResetRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	reason := c.FormValue("reason")

	if err := h.passwordResetService.RejectResetPassword(requestId, reason); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=rejection_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	if err := utils.SetFlashMessage(c, "success", "Reject success"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?success=reject_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
}
