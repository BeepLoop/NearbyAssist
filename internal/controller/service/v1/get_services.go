package service

import (
	"nearbyassist/internal/db/query/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetServices(c echo.Context) error {
	locations, err := query.GetServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, locations)
}
