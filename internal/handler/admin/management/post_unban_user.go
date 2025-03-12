package management

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) UnbanUser(c echo.Context) error {
	userId := c.Param("userId")

	if err := h.managementService.UnbanUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=unban_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
	}

	if err := utils.SetFlashMessage(c, "success", "unbanned user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?success=unban_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
