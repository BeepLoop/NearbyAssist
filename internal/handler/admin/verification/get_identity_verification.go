package verification

import (
	"context"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/identity_verification"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) GetRequestList(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	requests, err := h.verificationService.GetRequestList()
	if err != nil {
		fmt.Println(err.Error())
		page := pages.IdentityVerificationList(*admin, make([]dto.VerificationRequest, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.IdentityVerificationList(*admin, requests)
	return page.Render(context.Background(), c.Response().Writer)
}
