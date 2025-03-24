package complaint

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetReportedUsers(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	params := c.QueryParams()
	flash, _, _ := utils.RetrieveFlashMessage(c)

	results := make([]*models.ReportedUserModel, 0)

	if params.Has("query") {
		// Perform a search
	} else {
		limit, _ := strconv.Atoi(params.Get("limit"))
		if limit == 0 {
			limit = DEFAULT_LIMIT
		}

		offset, _ := strconv.Atoi(params.Get("offset"))
		if offset == 0 {
			offset = DEFAULT_OFFSET
		}

		users, err := h.complaintService.GetReportedUsers(limit, offset)
		if err != nil {
			empty := make([]models.ReportedUserModel, 0)
			page := pages.ReportedUsers(*admin, empty, flash)
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = users
	}

	data := make([]models.ReportedUserModel, 0)
	for _, report := range results {
		data = append(data, models.ReportedUserModel{
			Model:  report.Model,
			UserId: report.UserId,
			Reason: report.Reason,
			Detail: report.Detail,
		})
	}

	page := pages.ReportedUsers(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
