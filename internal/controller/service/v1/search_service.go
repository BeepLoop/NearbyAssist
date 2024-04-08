package service

import (
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SearchService(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	model := models.NewServiceModelWithLocation(params.Latitude, params.Longitude)

	services, err := model.GeoSpatialSearch(params.Query, params.Radius)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}
