package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func UnrestrictVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing vendor ID")
	}

	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	err = vendor_query.UnrestrictVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": id,
	})
}
