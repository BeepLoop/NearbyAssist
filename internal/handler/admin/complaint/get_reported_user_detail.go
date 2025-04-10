package complaint

import (
	"context"
	"nearbyassist/internal/dto"
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

	data, err := h.complaintService.GetReportedUserDetail(reportId)
	if err != nil {
		page := pages.ViewReportedUserDetail(*admin, dto.UserReport{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.ViewReportedUserDetail(*admin, *data)
	return page.Render(context.Background(), c.Response().Writer)
}
