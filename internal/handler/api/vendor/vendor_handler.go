package vendor

import (
	"nearbyassist/internal/models"
	vendor_service "nearbyassist/internal/service/vendor"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type vendorHandler struct {
	vendorService *vendor_service.Service
}

func NewHandler(vendorService *vendor_service.Service) *vendorHandler {
	return &vendorHandler{vendorService: vendorService}
}

func (h *vendorHandler) GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	vendor, err := h.vendorService.GetVendor(vendorId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "Vendor not found",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error while retrieving vendor",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendor,
	})
}

func (h *vendorHandler) GetVendorServiceList(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	services, err := h.vendorService.GetVendorServiceList(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor services",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}
