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
	email := c.FormValue("email")
	duration := c.FormValue("duration")

	if err := h.inviteService.Invite(username, email, duration); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
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
