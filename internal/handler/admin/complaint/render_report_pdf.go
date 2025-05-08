package complaint

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/views/pdftemplate"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) RenderReportPDF(c echo.Context) error {
	reportId := c.Param("reportId")

	data, err := h.complaintService.GetReport(reportId)
	if err != nil {
		page := pdftemplate.ReportUserTemplate(dto.UserReportDetail{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pdftemplate.ReportUserTemplate(*data)
	return page.Render(context.Background(), c.Response().Writer)
}
