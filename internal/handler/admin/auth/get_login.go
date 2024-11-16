package auth

import (
	"context"
	"nearbyassist/views/pages/auth"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *authHandler) GetLogin(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/?error=session_error")
	}

	_, ok := sess.Values["user"]
	if ok {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	page := pages.Login()
	return page.Render(context.Background(), c.Response().Writer)
}
