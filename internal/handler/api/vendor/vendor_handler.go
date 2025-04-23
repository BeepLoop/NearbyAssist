package vendor

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	resource_service "nearbyassist/internal/service/resource"
	vendor_service "nearbyassist/internal/service/vendor"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"
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

	response := response.Vendor{
		Id:       vendor.VendorId,
		Name:     vendor.User.Name,
		Email:    vendor.User.Email,
		ImageUrl: vendor.User.ImageUrl,
		Phone:    vendor.User.Phone,
		Rating:   vendor.Rating,
		Socials:  vendor.User.Socials,
		Expertise: slices.AppendSeq(
			make([]string, 0),
			utils.Map(vendor.Expertise, func(e models.ExpertiseModel) string { return e.Title }),
		),
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

	services, err := h.vendorService.GetVendorServicesList(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving vendor services",
			Error:   err.Error(),
		})
	}

	response := response.VendorServices{
		Vendor: response.Vendor{
			Id:       vendor.VendorId,
			Name:     vendor.User.Name,
			Email:    vendor.User.Email,
			ImageUrl: vendor.User.ImageUrl,
			Phone:    vendor.User.Phone,
			Rating:   vendor.Rating,
			Socials:  vendor.User.Socials,
			Expertise: slices.AppendSeq(
				make([]string, 0),
				utils.Map(vendor.Expertise, func(e models.ExpertiseModel) string { return e.Title }),
			),
		},
		Servics: slices.AppendSeq(
			make([]response.Service, 0),
			utils.Map(services, func(s *models.ServiceModel) response.Service {
				return response.Service{
					Id:          s.Id,
					VendorId:    s.VendorId,
					Title:       s.Title,
					Description: s.Description,
					Rate:        s.Rate,
					Tags: slices.AppendSeq(
						make([]response.Tag, 0),
						utils.Map(s.Tags, func(t *models.TagModel) response.Tag {
							return response.Tag{Id: t.Id, Title: t.Title}
						}),
					),
					Extras: slices.AppendSeq(
						make([]response.Extra, 0),
						utils.Map(s.Extras, func(x *models.ExtraModel) response.Extra {
							return response.Extra{
								Id:          x.Id,
								Title:       x.Title,
								Description: x.Description,
								Price:       x.Price,
							}
						}),
					),
					Images: slices.AppendSeq(
						make([]response.Image, 0),
						utils.Map(s.Images, func(i *models.ServicePhotoModel) response.Image {
							return response.Image{Id: i.Id, Url: i.Url}
						}),
					),
					Location: response.Location{
						Latitude:  s.Address.Latitude,
						Longitude: s.Address.Longitude,
					},
					Disabled: s.Disabled,
				}
			}),
		),
	}

	return c.JSON(http.StatusOK, response)
}
