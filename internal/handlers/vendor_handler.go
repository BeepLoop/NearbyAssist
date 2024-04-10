package handlers

import (
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/server"
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

func (h *vendorHandler) HandleCount(c echo.Context) error {
	filter := c.QueryParam("filter")

	model := models.NewVendorModel(h.server.DB)

	count, err := model.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"vendorCount": count,
	})
}

func (h *vendorHandler) HandleGetVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewVendorModel(h.server.DB)

	vendor, err := model.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	return c.JSON(http.StatusOK, vendor)
}

func (h *vendorHandler) HandleRestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	model := models.NewVendorModel(h.server.DB)

	err = model.RestrictAccount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": id,
	})
}

func (h *vendorHandler) HandleUnrestrict(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	model := models.NewVendorModel(h.server.DB)

	err = model.UnrestrictAccount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": id,
	})
}
