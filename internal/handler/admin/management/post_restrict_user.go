package management

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) RestrictUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/account-manangement/"+userId+"?error=Invalid_request")
	}

	if err := h.managementService.RestrictUser(userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-manangement/"+userId+"?error=error_restricting")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-manangement"+userId)
}
