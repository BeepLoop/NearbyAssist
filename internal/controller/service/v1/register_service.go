package service

import (
	"nearbyassist/internal/db/query/service"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterService(c echo.Context) error {
	var service types.ServiceRegister
	err := c.Bind(&service)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	err = c.Validate(service)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// construct point from given Latitude and Longitude
	transfromedData, err := utils.TransformServiceData(service)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = service_query.RegisterService(transfromedData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
