package auth

import (
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *authHandler) PostLogout(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

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

	_ = utils.SetFlashMessage(c, "success", "Logged out")

	activity := activitylog.Input{
		AdminId: admin.Id,
		Action:  activitylog.ACTION_LOGOUT,
	}
	activitylog.MustGetInstance().Create(activity)

	return c.Redirect(http.StatusSeeOther, "/auth/login")
}
