package service

import (
	review_query "nearbyassist/internal/db/query/review"
	service_query "nearbyassist/internal/db/query/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetServiceDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "serviceId is required",
		})
	}
	id, err := strconv.Atoi(strings.ReplaceAll(serviceId, "/", ""))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "service ID must be a number",
		})
	}

	serviceDetails, err := service_query.GetServiceDetails(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	reviewCount, err := review_query.ReviewCount(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	serviceDetails.ReviewCount = reviewCount

	servicePhotos, err := service_query.GetPhotos(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	serviceDetails.Photos = servicePhotos

	return c.JSON(http.StatusOK, serviceDetails)
}
