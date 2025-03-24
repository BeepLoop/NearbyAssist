package complaint

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetBugReports(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = DEFAULT_LIMIT
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = DEFAULT_OFFSET
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	complaints, err := h.complaintService.GetBugReports(limit, offset)
	if err != nil {
		page := pages.BugReports(*admin, make([]models.BugReportModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.BugReportModel, 0)
	for _, complaint := range complaints {
		date := utils.FormatDate(complaint.CreatedAt)

		images := make([]string, 0)
		for _, image := range complaint.Images {
			signedURL, err := h.resourceService.SignURLWithDefaultDuration(image)
			if err != nil {
				fmt.Println("Error generating signed url for bug image: ", err.Error())
				continue
			}

			images = append(images, signedURL)
		}

		data = append(data, models.BugReportModel{
			Id:          complaint.Id,
			Title:       complaint.Title,
			Detail:      complaint.Detail,
			CreatedAt:   date,
			CompletedAt: complaint.CompletedAt,
			Images:      images,
		})

	}

	page := pages.BugReports(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
