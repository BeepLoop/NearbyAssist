package middleware

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/login?error=unkown_session")
		}

		user := sess.Values["user"]
		_, ok := user.(models.AdminModel)
		if !ok {
			if err := utils.SetFlashMessage(c, "error", "invalid session"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_not_found")
			}

			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		return next(c)
	}
}
