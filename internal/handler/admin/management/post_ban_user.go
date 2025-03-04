package management

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) BanUser(c echo.Context) error {
	userId := c.Param("userId")

	if err := h.managementService.BanUser(userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=error_banning")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
