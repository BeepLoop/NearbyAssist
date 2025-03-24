package application

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/vendor_application"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplication(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	applications, err := h.applicationService.GetApplications()
	if err != nil {
		page := pages.Applications(*admin, make([]models.ApplicationModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.ApplicationModel, 0)
	for _, application := range applications {
		date := utils.FormatDate(application.CreatedAt)

		data = append(data, models.ApplicationModel{
			Model:                 models.Model{Id: application.Id, CreatedAt: date},
			UpdateableModel:       application.UpdateableModel,
			GeoSpatialModel:       application.GeoSpatialModel,
			ApplicantId:           application.ApplicantId,
			ExpertiseId:           application.ApplicantId,
			ApplicantName:         application.ApplicantName,
			SupportingDocumentUrl: application.SupportingDocumentUrl,
			PoliceClearanceUrl:    application.PoliceClearanceUrl,
			Status:                application.Status,
		})
	}

	page := pages.Applications(*admin, data)
	return page.Render(context.Background(), c.Response().Writer)
}
