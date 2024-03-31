package user

import (
	user_query "nearbyassist/internal/db/query/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountUsers(c echo.Context) error {
	users, err := user_query.CountUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
	})
}
