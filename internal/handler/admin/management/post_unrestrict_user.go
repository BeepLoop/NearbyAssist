package management

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) UnrestrictUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=Invalid_request")
	}

	if err := h.managementService.UnrestrictUser(userId); err != nil {
		if err := utils.SetFlashMessage(c, "error", err.Error()); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId+"?error=unrestrict_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/"+userId)
}
