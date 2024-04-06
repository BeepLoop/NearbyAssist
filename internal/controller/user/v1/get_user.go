package user

import (
	// "nearbyassist/internal/db/query/user"
	// "nearbyassist/internal/utils"
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	userModel := models.NewUserModel()
	user, err := userModel.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
