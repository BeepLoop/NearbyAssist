package service

import (
	"nearbyassist/internal/db/query/service"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetServices(c echo.Context) error {
	results, err := query.GetServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var services []types.TransformedServiceData
	for _, service := range results {
		lat, long := utils.LatlongExtractor(service.Location)
		services = append(services, types.TransformedServiceData{
			Vendor:      service.Vendor,
			Title:       service.Title,
			Description: service.Description,
			Rate:        service.Rate,
			Latitude:    lat,
			Longitude:   long,
		})
	}

	return c.JSON(http.StatusOK, services)
}
