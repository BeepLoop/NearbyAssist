package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) RestrictUser(c echo.Context) error {
	reason := c.FormValue("reason")
	duration := c.FormValue("duration")
	redirectRoute := c.FormValue("redirectRoute")
	password := c.FormValue("password")
	userId := c.Param("userId")

	admin, _ := utils.GetAdminFromSession(c)
	if err := h.managementService.RestrictUser(admin.Id, password, userId, reason, duration); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, redirectRoute+"?error=restrict_error")
		}

		return c.Redirect(http.StatusSeeOther, redirectRoute)
	}

	if err := utils.SetFlashMessage(c, "success", "restricted user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectRoute+"?success=restrict_success")
	}

	return c.Redirect(http.StatusSeeOther, redirectRoute)
}
