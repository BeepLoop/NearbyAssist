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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "vendorId is required",
		})
	}
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "vendor ID must be a number",
		})
	}

	vendor, err := user_query.GetVendor(id)
	if err != nil {
		return c.JSON(http.StatusNoContent, map[string]string{
			"error": err.Error(),
		})
	}

	reviewCount, err := review_query.ReviewCount(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	vendor.ReviewCount = reviewCount

	return c.JSON(http.StatusOK, vendor)
}
