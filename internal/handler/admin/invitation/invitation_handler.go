package invitation

import (
	"nearbyassist/internal/service/invite_service"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type handler struct {
	inviteService *invite_service.Service
}

func NewHandler(inviteService *invite_service.Service) *handler {
	return &handler{
		inviteService: inviteService,
	}
}

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
		if strings.Contains(err.Error(), "duplicate username") {
			if err := utils.SetFlashMessage(c, "error", "Username already in use"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=invite_error")
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

func (h *handler) JoinInvite(c echo.Context) error {
	code := c.QueryParam("code")

	if err := h.inviteService.Join(code); err != nil {
		if strings.Contains(err.Error(), "invitation expired") {
			if err := utils.SetFlashMessage(c, "error", "Invitation expired"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/?error=invitation_expired")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Error joining with invite link"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/?error=join_invite_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
