package vendor

import (
	"nearbyassist/internal/models"
	resource_service "nearbyassist/internal/service/resource"
	vendor_service "nearbyassist/internal/service/vendor"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type vendorHandler struct {
	vendorService   *vendor_service.Service
	resourceService *resource_service.Service
}

func NewHandler(vendorService *vendor_service.Service, resourceService *resource_service.Service) *vendorHandler {
	return &vendorHandler{
		vendorService:   vendorService,
		resourceService: resourceService,
	}
}

func (h *vendorHandler) GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	vendor, err := h.vendorService.FindById(vendorId)
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

	return c.JSON(http.StatusOK, vendor)
}

func (h *vendorHandler) GetVendorServiceList(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	_, err := h.vendorService.FindById(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor information",
			Error:   err.Error(),
		})
	}

	detail, err := h.vendorService.GetVendorServicesList(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor services",
			Error:   err.Error(),
		})
	}

	for _, service := range detail.Services {
		for _, image := range service.Images {
			signedURL, err := h.resourceService.SignURLWithDefaultDuration(image.Url)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error retrieving vendor services",
					Error:   err.Error(),
				})
			}
			image.Url = signedURL
		}
	}

	return c.JSON(http.StatusOK, detail)
}
