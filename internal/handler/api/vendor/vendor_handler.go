package vendor

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	vendor_service "nearbyassist/internal/service/vendor"
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

	response := struct {
		Id           string   `json:"id"`
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		Phone        string   `json:"phone"`
		ImageUrl     string   `json:"imageUrl"`
		Rating       string   `json:"rating"`
		IsRestricted bool     `json:"isRestricted"`
		Socials      []string `json:"socials"`
	}{
		Id:           vendor.Id,
		Name:         vendor.Name,
		Email:        vendor.Email,
		Phone:        vendor.Phone.String,
		ImageUrl:     vendor.ImageUrl,
		Rating:       vendor.Rating,
		IsRestricted: vendor.Restricted,
		Socials:      vendor.Socials,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *vendorHandler) GetVendorServiceList(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	vendor, err := h.vendorService.FindById(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor information",
			Error:   err.Error(),
		})
	}

	services, err := h.vendorService.GetVendorServiceList(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor services",
			Error:   err.Error(),
		})
	}

	response := response.VendorServiceList{
		Vendor: struct {
			Id           string   `json:"id"`
			Name         string   `json:"name"`
			Email        string   `json:"email"`
			Phone        string   `json:"phone"`
			ImageUrl     string   `json:"imageUrl"`
			Rating       string   `json:"rating"`
			IsRestricted bool     `json:"isRestricted"`
			Expertise    []string `json:"expertise"`
			Socials      []string `json:"socials"`
		}{
			Id:           vendor.Id,
			Name:         vendor.Name,
			Email:        vendor.Email,
			Phone:        vendor.Phone.String,
			ImageUrl:     vendor.ImageUrl,
			Rating:       vendor.Rating,
			IsRestricted: vendor.Restricted,
			Expertise:    vendor.Expertise,
			Socials:      vendor.Socials,
		},
	}

	for _, service := range services {
		response.Services = append(response.Services, struct {
			Id          string             `json:"id"`
			Title       string             `json:"title"`
			Description string             `json:"description"`
			Price       string             `json:"price"`
			Latitude    float64            `json:"latitude"`
			Longitude   float64            `json:"longitude"`
			Tags        []*models.TagModel `json:"tags"`
		}{
			Id:          service.Id,
			Title:       service.Title,
			Description: service.Description,
			Price:       service.Rate,
			Latitude:    service.Latitude,
			Longitude:   service.Longitude,
			Tags:        service.Tags,
		})
	}

	return c.JSON(http.StatusOK, response)
}
