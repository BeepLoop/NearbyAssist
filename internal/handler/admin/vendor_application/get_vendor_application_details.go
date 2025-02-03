package application

import (
	"context"
	"nearbyassist/internal/models"
	pages "nearbyassist/views/pages/vendor_application"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplicationDetails(c echo.Context) error {
	applicationId := c.Param("applicationId")

	application, err := h.applicationService.GetApplicationDetail(applicationId)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	supportingDocBase64, err := h.applicationService.GetFile(application.SupportingDocumentUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	policeClearanceBase64, err := h.applicationService.GetFile(application.PoliceClearanceUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := models.ApplicationModel{
		Model:                 application.Model,
		UpdateableModel:       application.UpdateableModel,
		GeoSpatialModel:       application.GeoSpatialModel,
		ApplicantId:           application.ApplicantId,
		ExpertiseId:           application.ExpertiseId,
		SupportingDocumentUrl: supportingDocBase64,
		PoliceClearanceUrl:    policeClearanceBase64,
		Status:                application.Status,
		ApplicantName:         application.ApplicantName,
		Expertise:             application.Expertise,
	}

	page := pages.VendorApplicationDetails(data)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *applicationHandler) AcceptRequest(c echo.Context) error {
	applicationId := c.Param("applicationId")
	_ = applicationId

	return nil
}

func (h *applicationHandler) RejectRequest(c echo.Context) error {
	applicationId := c.Param("applicationId")
	_ = applicationId

	return nil
}
