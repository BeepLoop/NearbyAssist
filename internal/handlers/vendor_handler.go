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
	status := models.VendorStatus(c.QueryParam("status"))

	count, err := h.server.DB.CountVendor(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	vendor, err := h.server.DB.FindVendorById(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "vendor not found")
	}

	// TODO: retrieve review count

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendor,
	})
}

func (h *vendorHandler) HandleRestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	if err := h.server.DB.RestrictVendor(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "vendor not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"restrictedId": vendorId,
	})
}

func (h *vendorHandler) HandleUnrestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	if err := h.server.DB.UnrestrictVendor(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "vendor not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"unrestrictedId": vendorId,
	})
}
