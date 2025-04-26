package passwordreset_handler

import (
	"context"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

type passwordResetHandler struct {
	passwordResetService *passwordreset_service.Service
}

func NewHandler(passwordResetService *passwordreset_service.Service) *passwordResetHandler {
	return &passwordResetHandler{
		passwordResetService: passwordResetService,
	}
}

func (h *passwordResetHandler) ChangePassword(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	page := pages.ChangePassword(flash)
	return page.Render(context.Background(), c.Response().Writer)
}
