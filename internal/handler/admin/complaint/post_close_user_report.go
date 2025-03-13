package complaint

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) CloseUserReport(c echo.Context) error {
	reportId := c.FormValue("reportId")
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	if err := h.complaintService.CloseUserReport(reportId, title, detail); err != nil {
		if err := utils.SetFlashMessage(c, "error", "failed to close report: "+reportId); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=report_close_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
	}

	if err := utils.SetFlashMessage(c, "success", "closed user report"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?success=closed_report")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
}
