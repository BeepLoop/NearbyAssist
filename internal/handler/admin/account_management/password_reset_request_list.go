package accountmanagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management/password_reset"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) PasswordResetRequestList(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	requests, err := h.passwordResetService.GetRequestList()
	if err != nil {
		page := pages.ResetRequests(*admin, make([]dto.PasswordResetRequest, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.ResetRequests(*admin, requests, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
