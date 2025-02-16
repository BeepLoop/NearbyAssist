package application

import (
	"context"
	"nearbyassist/internal/models"
	pages "nearbyassist/views/pages/vendor_application"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplication(c echo.Context) error {
	applications, err := h.applicationService.GetApplications()
	if err != nil {
		page := pages.Applications(make([]models.ApplicationModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.ApplicationModel, 0)
	for _, application := range applications {
		var date string

		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, application.CreatedAt)
		if err != nil {
			date = application.CreatedAt
		} else {
			date = t.Format(time.RFC1123)
		}

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

	page := pages.Applications(data)
	return page.Render(context.Background(), c.Response().Writer)
}
