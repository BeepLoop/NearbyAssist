package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	service_service "nearbyassist/internal/service/service"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type serviceHandler struct {
	service *service_service.Service
}

func NewHandler(service *service_service.Service) *serviceHandler {
	return &serviceHandler{service: service}
}

func (h *serviceHandler) CreateService(c echo.Context) error {
	req := new(request.NewServicePayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	serviceId, err := h.service.CreateService(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"service": serviceId,
	})
}

func (h *serviceHandler) GetService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	data, err := h.service.GetService(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error while retrieving service information",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"serviceInfo":    data["serviceInfo"],
		"vendorInfo":     data["vendorInfo"],
		"serviceImages":  data["serviceImages"],
		"countPerRating": data["countPerRating"],
	})
}

func (h *serviceHandler) UpdateService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	req := new(request.UpdateServicePayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   err.Error(),
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	if err := h.service.UpdateService(bearerToken, serviceId, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error updating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) SearchService(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())

	services, err := h.service.SearchService(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred while searching service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) FindServiceRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	origin := c.QueryParam("origin")

	polyline, err := h.service.FindRoute(serviceId, origin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error finding route to given destination",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"polyline": polyline,
	})
}
