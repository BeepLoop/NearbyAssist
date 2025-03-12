package management

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) BanUser(c echo.Context) error {
	userId := c.Param("userId")

	if err := h.managementService.BanUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=ban_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
