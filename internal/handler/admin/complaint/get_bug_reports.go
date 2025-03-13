package complaint

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"strconv"

	"github.com/labstack/echo/v4"
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

	flash, _, _ := utils.RetrieveFlashMessage(c)

	complaints, err := h.complaintService.GetBugReports(limit, offset)
	if err != nil {
		page := pages.BugReports(make([]models.BugReportModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.BugReportModel, 0)
	for _, complaint := range complaints {
		date := utils.FormatDate(complaint.CreatedAt)

		images := make([]string, 0)
		for _, image := range complaint.Images {
			base64Image, err := h.complaintService.GetFile(image)
			if err != nil {
				fmt.Println("error retrieving bug report image: ", err.Error())
				continue
			}

			images = append(images, base64Image)
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

	page := pages.BugReports(data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
