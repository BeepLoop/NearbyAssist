package handlers

import (
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type serviceHandler struct {
	server *server.Server
}

func NewServiceHandler(server *server.Server) *serviceHandler {
	return &serviceHandler{
		server: server,
	}
}

func (h *serviceHandler) HandleGetServices(c echo.Context) error {
	model := models.NewServiceModel()

	services, err := model.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}

func (h *serviceHandler) HandleRegisterService(c echo.Context) error {
	model := models.NewServiceModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	insertId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"serviceId": insertId,
	})
}

func (h *serviceHandler) HandleSearchService(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	model := models.NewServiceModelWithLocation(params.Latitude, params.Longitude)

	services, err := model.GeoSpatialSearch(params.Query, params.Radius)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}

func (h *serviceHandler) HandleGetDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	model := models.NewServiceModel()

	service, err := model.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	// TODO: retrieve photos associated with a given service

	return c.JSON(http.StatusOK, service)
}

func (h *serviceHandler) HandleGetByOwner(c echo.Context) error {
	ownerId := c.Param("ownerId")
	id, err := strconv.Atoi(ownerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	model := models.NewServiceModel()

	services, err := model.FindByVendorId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}
