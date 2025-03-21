package application

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/vendor_application"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplicationDetails(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	applicationId := c.Param("applicationId")

	application, err := h.applicationService.GetApplicationDetail(applicationId)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	supportingDocumentImage, err := h.resourceService.SignURLWithDefaultDuration(application.SupportingDocumentUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	policeClearanceImage, err := h.resourceService.SignURLWithDefaultDuration(application.PoliceClearanceUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(models.ApplicationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := models.ApplicationModel{
		Model: models.Model{
			Id:        application.Id,
			CreatedAt: utils.FormatDate(application.CreatedAt),
		},
		UpdateableModel:       application.UpdateableModel,
		GeoSpatialModel:       application.GeoSpatialModel,
		ApplicantId:           application.ApplicantId,
		ExpertiseId:           application.ExpertiseId,
		SupportingDocumentUrl: supportingDocumentImage,
		PoliceClearanceUrl:    policeClearanceImage,
		Status:                application.Status,
		ApplicantName:         application.ApplicantName,
		Expertise:             application.Expertise,
	}

	page := pages.VendorApplicationDetails(data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *applicationHandler) AcceptRequest(c echo.Context) error {
	applicationId := c.FormValue("applicationId")

	if applicationId == "" {
		if err := utils.SetFlashMessage(c, "error", "Invalid application ID"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=id_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	if err := h.applicationService.AcceptRequest(applicationId); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Request accept failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=accept_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}

func (h *applicationHandler) RejectRequest(c echo.Context) error {
	reason := c.FormValue("reason")
	applicationId := c.FormValue("applicationId")

	if applicationId == "" {
		if err := utils.SetFlashMessage(c, "error", "Invalid application ID"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=id_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	if err := h.applicationService.RejectRequest(applicationId, reason); err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "invalid reason") {
			if err := utils.SetFlashMessage(c, "error", "Provide a reason for rejection"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=invalid_reason_error")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=rejection_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}
