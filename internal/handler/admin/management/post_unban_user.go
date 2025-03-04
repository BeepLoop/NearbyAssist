package management

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) UnbanUser(c echo.Context) error {
	userId := c.Param("userId")

	if err := h.managementService.UnbanUser(userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=error_unbanning")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
