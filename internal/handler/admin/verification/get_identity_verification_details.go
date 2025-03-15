package verification

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/views/pages/identity_verification"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetIdentityVerificationDetails(c echo.Context) error {
	requestId := c.Param("requestId")

	request, err := h.verificationService.GetRequest(requestId)
	if err != nil {
		page := pages.IdentityVerificationDetails(models.IdentityVerificationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	frontIdBase64, err := h.resourceService.GetBase64File(request.FrontIdImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(models.IdentityVerificationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	backIdBase64, err := h.resourceService.GetBase64File(request.BackIdImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(models.IdentityVerificationModel{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	selfieBase64, err := h.resourceService.GetBase64File(request.FaceImageUrl)
	if err != nil {
		page := pages.IdentityVerificationDetails(models.IdentityVerificationModel{})
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
		FrontIdImageUrl: frontIdBase64,
		BackIdImageUrl:  backIdBase64,
		FaceImageUrl:    selfieBase64,
	}

	page := pages.IdentityVerificationDetails(data)
	return page.Render(context.Background(), c.Response().Writer)
}
