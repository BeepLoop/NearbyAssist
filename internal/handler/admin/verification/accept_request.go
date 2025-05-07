package verification

import (
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) AcceptRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")

	if requestId == "" {
		if err := utils.SetFlashMessage(c, "error", "Invalid application ID"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=id_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	user, _ := h.verificationService.GetUserWithRequest(requestId)

	if err := h.verificationService.AcceptRequest(requestId); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Approval failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=Failed_to_accept_request")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	admin, _ := utils.GetAdminFromSession(c)

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
