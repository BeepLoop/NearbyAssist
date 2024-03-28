package user

import (
	"nearbyassist/internal/db/query/user"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing user ID")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	user, err := user_query.GetUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
