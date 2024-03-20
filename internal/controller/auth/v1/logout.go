package auth

import (
	session_query "nearbyassist/internal/db/query/session"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandlLogout(c echo.Context) error {
	u := new(types.User)
	err := c.Bind(u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request data",
		})
	}

	err = session_query.LogoutSession(u.Name, u.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "unable to logout",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "logout successfull",
	})
}
