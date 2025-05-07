package accountmanagement

import (
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) RejectResetRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	reason := c.FormValue("reason")

	target, _ := h.passwordResetService.GetAdminFromRequestId(requestId)

	if err := h.passwordResetService.RejectResetPassword(requestId, reason); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=rejection_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	if err := utils.SetFlashMessage(c, "success", "Reject success"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?success=reject_success")
	}

	admin, _ := utils.GetAdminFromSession(c)

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_REJECTED_PASSWORD_RESET,
		TargetType: "admin",
		TargetId:   target.Id,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreasePasswordResetRequest()

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
}
