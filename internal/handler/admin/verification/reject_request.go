package verification

import (
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) RejectRequest(c echo.Context) error {
	reason := c.FormValue("reason")
	requestId := c.FormValue("requestId")

	if requestId == "" {
		if err := utils.SetFlashMessage(c, "error", "Invalid application ID"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=id_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	if err := h.verificationService.RejectRequest(requestId, reason); err != nil {
		if strings.Contains(err.Error(), "invalid reason") {
			if err := utils.SetFlashMessage(c, "error", "Invalid reason"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=invalid_reason_error")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=rejection_error")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/verification-requests")
}
