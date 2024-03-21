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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user id required",
		})
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user id must be a number",
		})
	}

	acquaintances, err := message_query.GetAcquaintances(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, acquaintances)
}
