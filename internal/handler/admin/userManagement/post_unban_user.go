package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) UnbanUser(c echo.Context) error {
	redirectRoute := c.FormValue("redirectRoute")
	password := c.FormValue("password")
	userId := c.Param("userId")

	admin, _ := utils.GetAdminFromSession(c)
	if err := h.managementService.UnbanUser(admin.Id, password, userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, redirectRoute+"?error=unban_error")
		}

		return c.Redirect(http.StatusSeeOther, redirectRoute)
	}

	if err := utils.SetFlashMessage(c, "success", "unbanned user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectRoute+"?success=unban_success")
	}

	return c.Redirect(http.StatusSeeOther, redirectRoute)
}
