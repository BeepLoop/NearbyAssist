package complaint

import (
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) CompleteBug(c echo.Context) error {
	bugId := c.FormValue("id")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.complaintService.CompleteBug(bugId); err != nil {
		if err := utils.SetFlashMessage(c, "error", "failed to mark bug as complete"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/bugs?error=completion_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/bugs")
	}

	if err := utils.SetFlashMessage(c, "success", "marked bug complete"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/complaints/bugs?success=completed_bug")
	}

	activity := activitylog.Input{
		AdminId: admin.Id,
		Action:  activitylog.ACTION_CLOSED_BUG_REPORT,
	}
	activitylog.MustGetInstance().Create(activity)

	sse.New().DecreaseBugReport()

	return c.Redirect(http.StatusSeeOther, "/admin/complaints/bugs")
}
