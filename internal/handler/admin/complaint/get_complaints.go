package complaint

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

func (h *complaintHandler) GetBugReports(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = DEFAULT_LIMIT
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = DEFAULT_OFFSET
	}

	complaints, err := h.complaintService.GetBugReports(limit, offset)
	if err != nil {
		page := pages.Complaints(make([]models.BugReportModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.BugReportModel, 0)
	for _, complaint := range complaints {
		date := utils.FormatDate(complaint.CreatedAt)

		data = append(data, models.BugReportModel{
			Model:           models.Model{Id: complaint.Id, CreatedAt: date},
			UpdateableModel: complaint.UpdateableModel,
			Title:           complaint.Title,
			Detail:          complaint.Detail,
		})
	}

	page := pages.Complaints(data)
	return page.Render(context.Background(), c.Response().Writer)
}
