package middleware

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := utils.GetAdminFromSession(c)
		if err != nil {
			if err := utils.SetFlashMessage(c, "error", "invalid session"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_session_error")
			}

			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		return next(c)
	}
}
