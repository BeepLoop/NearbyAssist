package service_vendor

import (
	review_query "nearbyassist/internal/db/query/review"
	"nearbyassist/internal/db/query/user"
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

	vendor, err := user_query.GetVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	reviewCount, err := review_query.ReviewCount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	vendor.ReviewCount = reviewCount

	return c.JSON(http.StatusOK, vendor)
}
