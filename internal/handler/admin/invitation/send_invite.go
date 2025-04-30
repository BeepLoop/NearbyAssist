package invitation

import (
	"nearbyassist/internal/service/invite_service"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) SendInvite(c echo.Context) error {
	username := c.FormValue("username")
	defaultPassword := c.FormValue("defaultPassword")
	email := c.FormValue("email")
	duration := c.FormValue("duration")

	if username == "" || email == "" || defaultPassword == "" {
		if err := utils.SetFlashMessage(c, "error", "Invalid invite information"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=invite_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts")
	}

	if err := h.inviteService.Invite(username, email, defaultPassword, duration); err != nil {
		if strings.Contains(err.Error(), invite_service.ERR_DUPLICATE_USERNAME) {
			if err := utils.SetFlashMessage(c, "error", invite_service.ERR_DUPLICATE_USERNAME); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=duplicate_username")
			}
		} else if strings.Contains(err.Error(), invite_service.ERR_DUPLICATE_EMAIL) {
			if err := utils.SetFlashMessage(c, "error", invite_service.ERR_DUPLICATE_EMAIL); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=duplicate_email")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Error sending invite"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=invite_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts")
	}

	if err := utils.SetFlashMessage(c, "success", "Invitation send"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?success=invitation_sent")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts")
}
