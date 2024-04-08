package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RestrictVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	model := models.NewVendorModel()

	err = model.RestrictAccount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": id,
	})
}
