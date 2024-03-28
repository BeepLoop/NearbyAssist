package service

import (
	"nearbyassist/internal/db/query/service"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SearchService(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := service_query.SearchServices(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	services := make([]types.TransformedServiceData, 0)
	for _, location := range result {
		lat, long := utils.LatlongExtractor(location.Location)
		services = append(services, types.TransformedServiceData{
			Id:          location.Id,
			Vendor:      location.Vendor,
			Title:       location.Title,
			Description: location.Description,
			Rate:        location.Rate,
			Latitude:    lat,
			Longitude:   long,
		})
	}

	return c.JSON(http.StatusOK, services)
}
