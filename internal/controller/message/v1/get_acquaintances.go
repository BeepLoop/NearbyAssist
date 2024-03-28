package message

import (
	message_query "nearbyassist/internal/db/query/message"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAcquaintances(c echo.Context) error {
	userId := c.QueryParam("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing user ID")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	acquaintances, err := message_query.GetAcquaintances(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, acquaintances)
}
