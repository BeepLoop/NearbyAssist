package management

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) UnrestrictUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/account-manangement/"+userId+"?error=Invalid_request")
	}

	if err := h.managementService.UnrestrictUser(userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-manangement/"+userId+"?error=error_unrestricting")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-manangement"+userId)
}
