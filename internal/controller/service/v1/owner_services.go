package service

import (
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetOwnerServices(c echo.Context) error {
	ownerId := c.Param("ownerId")
	id, err := strconv.Atoi(ownerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	model := models.NewServiceModel()

	services, err := model.FindByVendorId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}
