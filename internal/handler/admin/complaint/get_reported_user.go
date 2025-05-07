package complaint

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetReport(c echo.Context) error {
	reportId := c.Param("reportId")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	data, err := h.complaintService.GetReport(reportId)
	if err != nil {
		page := pages.ReportedUser(*admin, dto.UserReportDetail{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_VIEWED_REPORT,
		TargetType: "other",
		TargetId:   reportId,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)

	page := pages.ReportedUser(*admin, *data)
	return page.Render(context.Background(), c.Response().Writer)
}
