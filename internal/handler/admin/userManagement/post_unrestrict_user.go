package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) UnrestrictUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?error=Invalid_request")
	}

	if err := h.managementService.UnrestrictUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?error=unrestrict_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
	}

	if err := utils.SetFlashMessage(c, "success", "lifted restriction on user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?success=unrestrict_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
}
