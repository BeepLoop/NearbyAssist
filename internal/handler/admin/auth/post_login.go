package auth

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *authHandler) PostLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	admin, err := h.adminService.Login(username, password)
	if err != nil {
		if err := utils.SetFlashMessage(c, "error", "invalid credentials"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/login?error=login_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_error")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	sess.Values["user"] = admin
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_error")
	}

	_ = utils.SetFlashMessage(c, "success", "Logged in")

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}
