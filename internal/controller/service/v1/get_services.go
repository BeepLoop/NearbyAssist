package service

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetServices(c echo.Context) error {
	model := models.NewServiceModel()

	services, err := model.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}
