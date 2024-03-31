package service_vendor

import (
	review_query "nearbyassist/internal/db/query/review"
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing vendor ID")
	}

	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	vendor, err := vendor_query.GetVendor(id)
	if err != nil {
		if utils.DetermineNoRowsError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "vendor not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	reviewCount, err := review_query.ReviewCount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	vendor.ReviewCount = reviewCount

	return c.JSON(http.StatusOK, vendor)
}
