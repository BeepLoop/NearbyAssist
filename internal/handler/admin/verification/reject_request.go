package verification

import (
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/service/sse"
	verification_service "nearbyassist/internal/service/verification"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) RejectRequest(c echo.Context) error {
	reason := c.FormValue("reason")
	requestId := c.FormValue("requestId")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	user, _ := h.verificationService.GetUserWithRequest(requestId)

	if err := h.verificationService.RejectRequest(admin.Id, password, requestId, reason); err != nil {
		if strings.Contains(err.Error(), verification_service.ERR_INVALID_REASON) {
			if err := utils.SetFlashMessage(c, "error", "Invalid reason"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=invalid_reason_error")
			}
		} else if strings.Contains(err.Error(), verification_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=rejection_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_REJECTED_VERIFICATION,
		TargetType: "user",
		TargetId:   user.Id,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreaseVerification()

	return c.Redirect(http.StatusSeeOther, "/admin/verification-requests")
}
