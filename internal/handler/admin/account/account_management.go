package account

import (
	admin_service "nearbyassist/internal/service/admin"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type accountHandler struct {
	adminService *admin_service.Service
}

func NewHandler(adminService *admin_service.Service) *accountHandler {
	return &accountHandler{
		adminService: adminService,
	}
}

func (h *accountHandler) RequestPasswordReset(c echo.Context) error {
	username := c.FormValue("username")

	if err := h.adminService.RequestPasswordReset(username); err != nil {
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
