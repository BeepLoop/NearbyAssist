package userManagement

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) BanUser(c echo.Context) error {
	redirectRoute := c.FormValue("redirectRoute")
	userId := c.Param("userId")

	if err := h.managementService.BanUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, redirectRoute+"?error=ban_error")
		}

		return c.Redirect(http.StatusSeeOther, redirectRoute)
	}

	if err := utils.SetFlashMessage(c, "success", "banned user: "+userId); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectRoute+"?success=ban_success")
	}

	return c.Redirect(http.StatusSeeOther, redirectRoute)
}
