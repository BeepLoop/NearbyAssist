package complaint

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetReportedUserDetail(c echo.Context) error {
	reportId := c.Param("reportId")

	detail, err := h.complaintService.GetReportedUserDetail(reportId)
	if err != nil {
		if err := utils.SetFlashMessage(c, "error", "report not found"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=user_not_found_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
	}

	images := make([]string, 0)
	for _, image := range detail.Images {
		base64Image, err := h.complaintService.GetFile(image)
		if err != nil {
			fmt.Println("error retrieving image: ", err.Error())
			continue
		}

		images = append(images, base64Image)
	}

	data := models.ReportedUserModel{
		Model: models.Model{
			Id:        detail.Id,
			CreatedAt: utils.FormatDate(detail.CreatedAt),
		},
		UserId: detail.UserId,
		Reason: detail.Reason,
		Detail: detail.Detail,
		Images: images,
		Name:   detail.Name,
	}

	page := pages.ViewReportedUserDetail(data)
	return page.Render(context.Background(), c.Response().Writer)
}
