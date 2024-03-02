package service_vendor

import (
	query "nearbyassist/internal/db/query/user"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "vendorId is required",
		})
	}
	id := strings.ReplaceAll(vendorId, "/", "")

	vendor, err := query.GetVendor(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, vendor)
}
