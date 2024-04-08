package service

import (
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetServiceDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	model := models.NewServiceModel()

	service, err := model.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	// TODO: retrieve photos associated with a given service

	return c.JSON(http.StatusOK, service)
}
