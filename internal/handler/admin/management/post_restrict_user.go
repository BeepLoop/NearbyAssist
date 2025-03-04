package management

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) RestrictUser(c echo.Context) error {
	reason := c.FormValue("reason")
	duration := c.FormValue("duration")
	userId := c.Param("userId")

	if err := h.managementService.RestrictUser(userId, reason, duration); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=error_restricting")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
