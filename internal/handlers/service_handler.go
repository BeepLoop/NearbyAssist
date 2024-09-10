package handlers

import (
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"
	"strings"

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
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	services, err := service.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting services",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) HandleCount(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	count, err := service.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting service count",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *serviceHandler) HandleRegisterService(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if err := c.Bind(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   err.Error(),
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	// Validate that the user is a registered vendor
	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := vendor.FindById(id); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "You are not a registered vendor",
			Error:   err.Error(),
		})
	}

	if _, err := service.EncryptDescription(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting description",
			Error:   encryption.ENCRYPTION_ERR,
		})
	}

	if _, err := service.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"serviceId": service.Id,
	})
}

func (h *serviceHandler) HandleGetDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	// Get service  info
	if _, err := service.FindById(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Service not found",
			Error:   err.Error(),
		})
	}

	if tags, err := service.GetTags(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting service tags",
			Error:   err.Error(),
		})
	} else {
		service.Tags = tags
	}

	// Get vendor info
	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := vendor.FindByServiceId(service.Id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Vendor not found",
			Error:   err.Error(),
		})
	}

	if _, err := vendor.DecryptVendor(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting vendor",
			Error:   encryption.DECRYPTION_ERR,
		})
	}

	// Get count per review rating
	reviews, err := service.GetReviews()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting reviews",
			Error:   err.Error(),
		})
	}

	countPerRating := response.NewCountPerRating()
	for _, review := range reviews {
		switch review.Rating {
		case 5:
			countPerRating["five"]++
		case 4:
			countPerRating["four"]++
		case 3:
			countPerRating["three"]++
		case 2:
			countPerRating["two"]++
		case 1:
			countPerRating["one"]++
		}
	}

	// Get service images
	photos, err := service.GetPhotos()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting service images",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"serviceInfo":    service,
		"vendorInfo":     vendor,
		"serviceImages":  photos,
		"countPerRating": countPerRating,
	})
}

func (h *serviceHandler) HandleUpdateService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if err := c.Bind(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   err.Error(),
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	// Validate if the service id  is owned by the requester
	if vendor, err := service.GetVendor(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting vendor",
			Error:   err.Error(),
		})
	} else {
		if vendor.VendorId != id {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You do not own this service",
				Error:   "You do not own this service",
			})
		}
	}

	if _, err := service.EncryptDescription(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting description",
			Error:   encryption.ENCRYPTION_ERR,
		})
	}

	if err := service.Update(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error updating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":   "Service updated",
		"serviceId": serviceId,
	})
}

func (h *serviceHandler) HandleDeleteService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	// Validate if the service is owned by the requester
	if vendor, err := service.GetVendor(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting vendor",
			Error:   err.Error(),
		})
	} else {
		if vendor.VendorId != id {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You do not own this service",
				Error:   "You do not own this service",
			})
		}
	}

	if err := service.Delete(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error deleting service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) HandleSearchService(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	params := utils.ParseQuery(c.QueryString())
	result, err := service.GeoSpatialSearch(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error searching services",
			Error:   err.Error(),
		})
	}

	// TODO: rank services by suggestability
	searchResult := make([]response.SearchResult, 0)
	var scoreError error
	for _, service := range result {
		score, err := h.server.SuggestionEngine.GenerateSuggestability(service)
		if err != nil {
			scoreError = err
			break
		}

		decrypted, err := h.server.Encrypt.DecryptString(service.Vendor)
		if err != nil {
			scoreError = err
			break
		}

		res := response.SearchResult{
			Id:             service.Id,
			Suggestability: score,
			Vendor:         decrypted,
			Latitude:       service.Latitude,
			Longitude:      service.Longitude,
		}

		searchResult = append(searchResult, res)
	}

	if scoreError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error generating suggestability score",
			Error:   scoreError.Error(),
		})
	}

	// sort services by Suggestability
	sortedResult := utils.BubbleSort(searchResult)

	// TODO: Generate ranks for the sorted services

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": sortedResult,
	})
}

func (h *serviceHandler) HandleGetByVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	services, err := service.FindByAllVendorId(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting services",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

// takes origin as QueryString ex: origin=lat,long
func (h *serviceHandler) HandleFindRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := service.FindById(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Service not found",
			Error:   err.Error(),
		})
	}

	origin, err := parseOrigin(c.QueryParam("origin"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid origin",
			Error:   err.Error(),
		})
	}

	destination := models.NewGeoGeoSpatialModelWithData(service.Latitude, service.Longitude)
	polyline, err := h.server.RouteEngine.FindRoute(origin, destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error finding route",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"polyline": polyline,
	})
}

func parseOrigin(query string) (*models.GeoSpatialModel, error) {
	coords := strings.Split(query, ",")
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, err
	}

	origin := models.NewGeoGeoSpatialModelWithData(lat, long)

	return origin, nil
}
