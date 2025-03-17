package passwordreset_handler

import (
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	"net/http"

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

func (h *passwordResetHandler) RequestPasswordReset(c echo.Context) error {
	username := c.FormValue("username")

	if err := h.passwordResetService.RequestPasswordReset(username); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Request failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/?error=request_error")
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	if err := utils.SetFlashMessage(c, "success", "Request submitted"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/?success=request_submitted")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
