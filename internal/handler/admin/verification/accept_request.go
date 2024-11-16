package verification

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *verificationHandler) AcceptRequest(c echo.Context) error {
	requestId := c.Param("requestId")
	if requestId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=Invalid_request")
	}

	if err := h.verificationService.AcceptRequest(requestId); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId+"?error=Failed_to_accept_request")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/verification-requests/"+requestId)
}
