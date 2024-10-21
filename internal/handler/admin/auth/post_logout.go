package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *authHandler) PostLogout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=session_error")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=session_error")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
