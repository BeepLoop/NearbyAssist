package service

import (
	review_query "nearbyassist/internal/db/query/review"
	service_query "nearbyassist/internal/db/query/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetServiceDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing service ID")
	}

	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	serviceDetails, err := service_query.GetServiceDetails(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	reviewCount, err := review_query.ReviewCount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	serviceDetails.ReviewCount = reviewCount

	servicePhotos, err := service_query.GetPhotos(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	serviceDetails.Photos = servicePhotos

	return c.JSON(http.StatusOK, serviceDetails)
}
