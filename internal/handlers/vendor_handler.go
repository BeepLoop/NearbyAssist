package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type vendorHandler struct {
	server *server.Server
}

func NewVendorHandler(server *server.Server) *vendorHandler {
	return &vendorHandler{server}
}

func (h *vendorHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Vendor base route",
	})
}

func (h *vendorHandler) HandleCount(c echo.Context) error {
	param := c.QueryParam("filter")
	var filter models.VendorStatusFilter
	switch param {
	case "":
		filter = models.VENDOR_STATUS_ALL
	case "all":
		filter = models.VENDOR_STATUS_ALL
	case "restricted":
		filter = models.VENDOR_STATUS_RESTRICTED
	case "unrestricted":
		filter = models.VENDOR_STATUS_UNRESTRICTED
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter found")
	}

	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	count, err := vendor.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving vendor count")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *vendorHandler) HandleGetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := vendor.FindById(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Vendor not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendor,
	})
}

func (h *vendorHandler) HandleRestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := vendor.Restrict(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while restricting vendor")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendorId,
	})
}

func (h *vendorHandler) HandleUnrestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := vendor.Unrestrict(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while unrestricting vendor")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendorId,
	})
}
