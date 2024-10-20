package middleware

import (
	"nearbyassist/internal/models"
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
			return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_not_found")
		}

		return next(c)
	}
}
