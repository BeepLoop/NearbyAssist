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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = session_query.LogoutSession(u.Name, u.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "logout successfull",
	})
}
