package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

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
		"vendorCount": count,
	})
}

func (h *vendorHandler) HandleGetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	vendor, err := h.server.DB.FindVendorById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	return c.JSON(http.StatusOK, utils.Mapper{
		"vendor": vendor,
	})
}

func (h *vendorHandler) HandleRestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	err = h.server.DB.RestrictVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"restricedId": id,
	})
}

func (h *vendorHandler) HandleUnrestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	err = h.server.DB.UnrestrictVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"unrestrictedId": id,
	})
}
