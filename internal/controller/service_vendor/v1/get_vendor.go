package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewVendorModel()

	vendor, err := model.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	return c.JSON(http.StatusOK, vendor)
}
