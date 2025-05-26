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

func (h *verificationHandler) AcceptRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	user, _ := h.verificationService.GetUserWithRequest(requestId)

	if err := h.verificationService.AcceptRequest(admin.Id, password, requestId); err != nil {
		if strings.Contains(err.Error(), verification_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=Unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Approval failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=Failed_to_accept_request")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_APPROVED_VERIFICATION,
		TargetType: "user",
		TargetId:   user.Id,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreaseVerification()

	return c.Redirect(http.StatusSeeOther, "/admin/verification-requests")
}
