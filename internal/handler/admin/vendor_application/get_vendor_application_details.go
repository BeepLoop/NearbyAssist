package application

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/activitylog"
	application_service "nearbyassist/internal/service/application"
	"nearbyassist/internal/service/cache"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/vendor_application"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplicationDetails(c echo.Context) error {
	params := c.QueryParams()

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	applicationId := c.Param("applicationId")

	var application *models.ApplicationModel

	if params.Has("fresh") && params.Get("fresh") == "true" {
		res, err := h.applicationService.GetApplicationDetail(applicationId)
		if err != nil {
			page := pages.VendorApplicationDetails(*admin, models.ApplicationModel{}, flash)
			return page.Render(context.Background(), c.Response().Writer)
		}

		application = res
		cache.NewGoCache().Set(c.Request().RequestURI, res)
	} else {
		inCache, exists := cache.NewGoCache().Get(c.Request().RequestURI)
		if exists {
			application = inCache.(*models.ApplicationModel)
		} else {
			res, err := h.applicationService.GetApplicationDetail(applicationId)
			if err != nil {
				page := pages.VendorApplicationDetails(*admin, models.ApplicationModel{}, flash)
				return page.Render(context.Background(), c.Response().Writer)
			}

			application = res
			cache.NewGoCache().Set(c.Request().RequestURI, res)
		}
	}

	supportingDocumentImage, err := h.resourceService.SignURLWithDefaultDuration(application.SupportingDocumentUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(*admin, models.ApplicationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	policeClearanceImage, err := h.resourceService.SignURLWithDefaultDuration(application.PoliceClearanceUrl)
	if err != nil {
		page := pages.VendorApplicationDetails(*admin, models.ApplicationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := models.ApplicationModel{
		Model: models.Model{
			Id:        application.Id,
			CreatedAt: utils.FormatDate(application.CreatedAt),
		},
		UpdatedAt:             application.UpdatedAt,
		GeoSpatialModel:       application.GeoSpatialModel,
		ApplicantId:           application.ApplicantId,
		ExpertiseId:           application.ExpertiseId,
		SupportingDocumentUrl: supportingDocumentImage,
		PoliceClearanceUrl:    policeClearanceImage,
		Status:                application.Status,
		ApplicantName:         application.ApplicantName,
		Expertise:             application.Expertise,
	}

	page := pages.VendorApplicationDetails(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *applicationHandler) AcceptRequest(c echo.Context) error {
	applicationId := c.FormValue("applicationId")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	application, _ := h.applicationService.GetApplicationDetail(applicationId)

	if err := h.applicationService.AcceptRequest(admin.Id, password, applicationId); err != nil {
		if strings.Contains(err.Error(), application_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Request accept failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=accept_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_APPROVED_APPLICATION,
		TargetType: "user",
		TargetId:   application.ApplicantId,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreaseApplication()

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}

func (h *applicationHandler) RejectRequest(c echo.Context) error {
	reason := c.FormValue("reason")
	applicationId := c.FormValue("applicationId")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	application, _ := h.applicationService.GetApplicationDetail(applicationId)

	if err := h.applicationService.RejectRequest(admin.Id, password, applicationId, reason); err != nil {
		if strings.Contains(err.Error(), "invalid reason") {
			if err := utils.SetFlashMessage(c, "error", "Provide a reason for rejection"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=invalid_reason_error")
			}
		} else if strings.Contains(err.Error(), application_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId+"?error=rejection_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications/"+applicationId)
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_REJECTED_APPLICATION,
		TargetType: "user",
		TargetId:   application.ApplicantId,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreaseApplication()

	return c.Redirect(http.StatusSeeOther, "/admin/vendor-applications")
}
