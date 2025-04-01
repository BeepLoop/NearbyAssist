package complaint

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/cache"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetReportedUserDetail(c echo.Context) error {
	params := c.QueryParams()

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	reportId := c.Param("reportId")

	var reportedUser *models.ReportedUserModel

	if params.Has("fresh") && params.Get("fresh") == "true" {
		detail, err := h.complaintService.GetReportedUserDetail(reportId)
		if err != nil {
			if err := utils.SetFlashMessage(c, "error", "report not found"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=user_not_found_error")
			}

			return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
		}

		reportedUser = detail
	} else {
		inCache, exists := cache.NewGoCache().Get(c.Request().RequestURI)
		if exists {
			reportedUser = inCache.(*models.ReportedUserModel)
		} else {
			detail, err := h.complaintService.GetReportedUserDetail(reportId)
			if err != nil {
				if err := utils.SetFlashMessage(c, "error", "report not found"); err != nil {
					return c.Redirect(http.StatusSeeOther, "/admin/complaints/users?error=user_not_found_error")
				}

				return c.Redirect(http.StatusSeeOther, "/admin/complaints/users")
			}

			reportedUser = detail
		}
	}

	images := make([]string, 0)
	for _, image := range reportedUser.Images {
		signedURL, err := h.resourceService.SignURLWithDefaultDuration(image)
		if err != nil {
			fmt.Println("Error generating signed url: ", err.Error())
			continue
		}

		images = append(images, signedURL)
	}

	data := models.ReportedUserModel{
		Model: models.Model{
			Id:        reportedUser.Id,
			CreatedAt: utils.FormatDate(reportedUser.CreatedAt),
		},
		ReportedBy: reportedUser.ReportedBy,
		UserId:     reportedUser.UserId,
		Reason:     reportedUser.Reason,
		Detail:     reportedUser.Detail,
		Images:     images,
		Name:       reportedUser.Name,
	}

	page := pages.ViewReportedUserDetail(*admin, data)
	return page.Render(context.Background(), c.Response().Writer)
}
