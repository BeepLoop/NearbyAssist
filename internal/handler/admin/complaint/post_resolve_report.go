package complaint

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) Resolve(c echo.Context) error {
	reportId := c.FormValue("reportId")
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	if err := h.complaintService.ActOnReport(reportId, title, detail, "resolved"); err != nil {
		if err := utils.SetFlashMessage(c, "error", "failed to resolve report: "+reportId); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=resolve_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
	}

	if err := utils.SetFlashMessage(c, "success", "resolved report"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?success=resolved_report")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
}
