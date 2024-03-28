package service

import (
	service_query "nearbyassist/internal/db/query/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetOwnerServices(c echo.Context) error {
	ownerId := c.Param("ownerId")
	if ownerId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing owner ID")
	}

	id, err := strconv.Atoi(ownerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	services, err := service_query.GetOwnerServices(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}
