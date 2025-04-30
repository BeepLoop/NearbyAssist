package invitation

import (
	"nearbyassist/internal/service/invite_service"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) JoinInvite(c echo.Context) error {
	code := c.QueryParam("code")

	if err := h.inviteService.Join(code); err != nil {
		if strings.Contains(err.Error(), invite_service.ERR_EXPIRED_INVITE) {
			if err := utils.SetFlashMessage(c, "error", invite_service.ERR_EXPIRED_INVITE); err != nil {
				return c.Redirect(http.StatusSeeOther, "/?error=invitation_expired")
			}
		} else if strings.Contains(err.Error(), invite_service.ERR_CODE_NOT_FOUND) {
			if err := utils.SetFlashMessage(c, "error", invite_service.ERR_CODE_NOT_FOUND); err != nil {
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
