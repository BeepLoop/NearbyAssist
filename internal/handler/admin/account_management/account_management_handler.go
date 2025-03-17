package accountmanagement

import (
	"context"
	"nearbyassist/internal/models"
	admin_service "nearbyassist/internal/service/admin"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management"

	"github.com/labstack/echo/v4"
)

type accountManagementHandler struct {
	adminService         *admin_service.Service
	passwordResetService *passwordreset_service.Service
}

func NewHandler(adminService *admin_service.Service, passwordResetService *passwordreset_service.Service) *accountManagementHandler {
	return &accountManagementHandler{
		adminService:         adminService,
		passwordResetService: passwordResetService,
	}
}

func (h accountManagementHandler) AddAccount(c echo.Context) error {
	page := pages.AddAccount()
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *accountManagementHandler) ResetRequests(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	requests, err := h.passwordResetService.GetResetRequests()
	if err != nil {
		c.Logger().Warnf("Error getting reset requests: %s", err.Error())
		page := pages.ResetRequests(make([]models.PasswordResetRequestModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.PasswordResetRequestModel, 0)
	for _, request := range requests {
		data = append(data, models.PasswordResetRequestModel{
			Model: models.Model{
				Id:        request.Id,
				CreatedAt: utils.FormatDate(request.CreatedAt),
			},
			AdminId:  request.AdminId,
			Username: request.Username,
		})
	}

	page := pages.ResetRequests(data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
