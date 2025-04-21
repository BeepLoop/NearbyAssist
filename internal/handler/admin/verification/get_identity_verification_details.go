package verification

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/identity_verification"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetRequest(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	data, err := h.verificationService.GetRequest(c.Param("requestId"))
	if err != nil {
		page := pages.IdentityVerification(*admin, dto.VerificationRequest{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.IdentityVerification(*admin, *data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
