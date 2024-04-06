package user

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountUsers(c echo.Context) error {
	model := models.NewUserModel()

	count, err := model.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userCount": count,
	})
}
