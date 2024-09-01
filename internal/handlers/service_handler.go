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
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	services, err := service.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) HandleCount(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	count, err := service.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while counting services")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *serviceHandler) HandleRegisterService(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := c.Bind(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error occurred binding request data")
	}

	if err := c.Validate(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	// Validate that the user is a registered vendor
	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := vendor.FindById(id); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "User is not a registered vendor")
	}

	if _, err := service.EncryptDescription(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := service.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while creating service")
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"serviceId": service.Id,
	})
}

func (h *serviceHandler) HandleGetDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	// Get service  info
	if _, err := service.FindById(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "service not found")
	}

	if tags, err := service.GetTags(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving service tags")
	} else {
		service.Tags = tags
	}

	// Get vendor info
	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := vendor.FindByServiceId(service.Id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Vendor of the service not found")
	}

	if _, err := vendor.DecryptVendor(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	// Get count per review rating
	reviews, err := service.GetReviews()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := c.Bind(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	if err := c.Validate(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	// Validate if the service id  is owned by the requester
	if vendor, err := service.GetVendor(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving service vendor")
	} else {
		if vendor.VendorId != id {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not own this service")
		}
	}

	if _, err := service.EncryptDescription(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if err := service.Update(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while updating service")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":   "Service updated",
		"serviceId": serviceId,
	})
}

func (h *serviceHandler) HandleDeleteService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	// Validate if the service is owned by the requester
	if vendor, err := service.GetVendor(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving service vendor")
	} else {
		if vendor.VendorId != id {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not own this service")
		}
	}

	if err := service.Delete(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while deleting service")
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) HandleSearchService(c echo.Context) error {
	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	params := utils.ParseQuery(c.QueryString())
	result, err := service.GeoSpatialSearch(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error occurred while searching services")
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
		return echo.NewHTTPError(http.StatusInternalServerError, scoreError.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	services, err := service.FindByAllVendorId(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving services")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

// takes origin as QueryString ex: origin=lat,long
func (h *serviceHandler) HandleFindRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service := models.NewServiceModel(h.server.IdGen, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := service.FindById(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "service not found")
	}

	origin, err := parseOrigin(c.QueryParam("origin"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid origin coordinates")
	}

	destination := models.NewGeoGeoSpatialModelWithData(service.Latitude, service.Longitude)
	polyline, err := h.server.RouteEngine.FindRoute(origin, destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find routes at the moment")
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
