package middleware

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EnsureStaff(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		activeSession, err := utils.GetAdminFromSession(c)
		if err != nil {
			if err := utils.SetFlashMessage(c, "error", "invalid session"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/auth/login?error=invalid_session_error")
			}

			return c.Redirect(http.StatusSeeOther, "/auth/login")
		}

		// NOTE: update redirect to home route if role is not admin
		if activeSession.Role != "staff" {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}
