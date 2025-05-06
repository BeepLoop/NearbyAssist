package complaint

import (
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) Close(c echo.Context) error {
	reportId := c.FormValue("reportId")
	action := c.FormValue("action")
	note := c.FormValue("note")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.complaintService.ActOnReport(reportId, action, admin.Id, note); err != nil {
		if err := utils.SetFlashMessage(c, "error", "failed to resolve report: "+reportId); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=resolve_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
	}

	if err := utils.SetFlashMessage(c, "success", "resolved report"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?success=resolved_report")
	}

	sse.New().Report--

	return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
}
