package verification

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/identity_verification"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetIdentityVerification(c echo.Context) error {
	requests, err := h.verificationService.GetIdentityVerificationRequests()
	if err != nil {
		page := pages.IdentityVerification(make([]models.IdentityVerificationModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.IdentityVerificationModel, 0)
	for _, request := range requests {
		data = append(data, models.IdentityVerificationModel{
			Model: models.Model{
				Id:        request.Id,
				CreatedAt: utils.FormatDate(request.CreatedAt),
			},
			UserId:          request.UserId,
			Status:          request.Status,
			Name:            request.Name,
			Address:         request.Address,
			IdType:          request.IdType,
			IdNumber:        request.IdNumber,
			FrontIdImageUrl: request.FrontIdImageUrl,
			BackIdImageUrl:  request.BackIdImageUrl,
			FaceImageUrl:    request.FaceImageUrl,
		})
	}

	page := pages.IdentityVerification(data)
	return page.Render(context.Background(), c.Response().Writer)
}
