package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) BanUser(c echo.Context) error {
	userId := c.Param("userId")

	if err := h.managementService.BanUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?error=ban_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
	}

	if err := utils.SetFlashMessage(c, "success", "banned user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId+"?success=ban_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/user-management/"+userId)
}
