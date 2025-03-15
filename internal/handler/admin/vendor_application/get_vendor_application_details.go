package application

import (
	"context"
	"nearbyassist/internal/models"
	pages "nearbyassist/views/pages/vendor_application"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplicationDetails(c echo.Context) error {
	applicationId := c.Param("applicationId")

	application, err := h.applicationService.GetApplicationDetail(applicationId)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	supportingDocBase64, err := h.resourceService.GetBase64File(application.SupportingDocumentUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	policeClearanceBase64, err := h.resourceService.GetBase64File(application.PoliceClearanceUrl)
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
	if applicationId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=Invalid_request")
	}

	if err := h.applicationService.AcceptRequest(applicationId); err != nil {
		return c.Redirect(
			http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=Failed_to_accept_request")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}

func (h *applicationHandler) RejectRequest(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=Invalid_request")
	}

	reason := "default reason"

	if err := h.applicationService.RejectRequest(applicationId, reason); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=Failed_to_accept_request")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}
