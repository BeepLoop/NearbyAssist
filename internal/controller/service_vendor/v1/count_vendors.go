package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountVendors(c echo.Context) error {
	filter := c.QueryParam("filter")

	model := models.NewVendorModel()

	count, err := model.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"vendorCount": count,
	})
}
