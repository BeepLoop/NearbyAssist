package user

import (
	"nearbyassist/internal/db/query/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(strings.ReplaceAll(userId, "/", ""))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	user, err := user_query.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
