package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) RestrictUser(c echo.Context) error {
	reason := c.FormValue("reason")
	duration := c.FormValue("duration")
	userId := c.Param("userId")

	if err := h.managementService.RestrictUser(userId, reason, duration); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?error=restrict_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
	}

	if err := utils.SetFlashMessage(c, "success", "restricted user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?success=restrict_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
}
