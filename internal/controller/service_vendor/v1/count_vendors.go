package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountVendors(c echo.Context) error {
	filter := c.QueryParam("filter")

	var vendorCount int
	switch filter {
	case "restricted":
		count, err := vendor_query.CountRestrictedVendors()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		vendorCount = count
	case "unrestricted":
		count, err := vendor_query.CountUnrestrictedVendors()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		vendorCount = count
	default:
		count, err := vendor_query.CountVendors()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		vendorCount = count
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"vendorCount": vendorCount,
	})
}
