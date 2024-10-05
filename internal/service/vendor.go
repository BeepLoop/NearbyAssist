package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/store/vendor"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type VendorService struct {
	store vendor.VendorStore
}

func NewVendorService(store vendor.VendorStore) *VendorService {
	return &VendorService{
		store: store,
	}
}

func (s *VendorService) GetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	vendor, err := s.store.FindById(vendorId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "Vendor not found",
				Error:   "Vendor not found",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error finding vendor",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendor,
	})
}
