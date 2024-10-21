package verification

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetIdentityVerification(c echo.Context) error {
	requests, err := h.verificationService.GetIdentityVerificationRequests()
	if err != nil {
		page := pages.IdentityVerification(make([]models.IdentityVerificationModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.IdentityVerification(requests)
	return page.Render(context.Background(), c.Response().Writer)
}
