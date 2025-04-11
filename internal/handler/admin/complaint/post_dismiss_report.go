package complaint

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) Dismiss(c echo.Context) error {
	reportId := c.FormValue("reportId")
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	if err := h.complaintService.ActOnReport(reportId, title, detail, "dismissed"); err != nil {
		if err := utils.SetFlashMessage(c, "error", "failed to dismiss report: "+reportId); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=dismiss_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
	}

	if err := utils.SetFlashMessage(c, "success", "dismissed report"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?success=dismissed_report")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
}
