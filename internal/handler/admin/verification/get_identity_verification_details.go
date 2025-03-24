package verification

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/identity_verification"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetIdentityVerificationDetails(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	requestId := c.Param("requestId")

	request, err := h.verificationService.GetRequest(requestId)
	if err != nil {
		page := pages.IdentityVerificationDetails(*admin, models.IdentityVerificationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	fontIdImage, err := h.resourceService.SignURLWithDefaultDuration(request.FrontIdImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(*admin, models.IdentityVerificationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	backIdImage, err := h.resourceService.SignURLWithDefaultDuration(request.BackIdImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(*admin, models.IdentityVerificationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	selfieImage, err := h.resourceService.SignURLWithDefaultDuration(request.FaceImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(*admin, models.IdentityVerificationModel{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := models.IdentityVerificationModel{
		Model:           request.Model,
		UpdateableModel: request.UpdateableModel,
		UserId:          request.UserId,
		Name:            request.Name,
		Address:         request.Address,
		IdType:          request.IdType,
		IdNumber:        request.IdNumber,
		Status:          request.Status,
		FrontIdImageUrl: fontIdImage,
		BackIdImageUrl:  backIdImage,
		FaceImageUrl:    selfieImage,
	}

	page := pages.IdentityVerificationDetails(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
