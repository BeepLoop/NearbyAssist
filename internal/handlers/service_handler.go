package handlers

import (
	"math/rand"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
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
	services, err := h.server.DB.FindAllService()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) HandleRegisterService(c echo.Context) error {
	req := &request.NewService{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.VendorId = userId
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate that the user is a registered vendor
	if _, err := h.server.DB.FindVendorById(req.VendorId); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "user is not a registered vendor")
	}

	insertId, err := h.server.DB.RegisterService(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"serviceId": insertId,
	})
}

func (h *serviceHandler) HandleUpdateService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	req := &request.UpdateService{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.VendorId = userId
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	} else {
		req.Id = id
	}

	// Validate if the service id  is owned by the requestor
	if owner, err := h.server.DB.FindServiceOwner(req.Id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		if owner.Id != req.VendorId {
			return echo.NewHTTPError(http.StatusForbidden, "you do not own this service")
		}
	}

	if err := h.server.DB.UpdateService(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":   "Service updated",
		"serviceId": serviceId,
	})
}

func (h *serviceHandler) HandleDeleteService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Validate if the service is owned by the requestor
	if owner, err := h.server.DB.FindServiceOwner(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		if owner.Id != userId {
			return echo.NewHTTPError(http.StatusForbidden, "you do not own this service")
		}
	}

	if err := h.server.DB.DeleteService(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) HandleSearchService(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	services, err := h.server.DB.GeoSpatialSearch(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: rank services by suggestability

	searchResult := make([]response.SearchResult, 0)
	for _, service := range services {
		res := response.SearchResult{
			Id:             service.Id,
			Suggestability: rand.Float32(),
			Vendor:         service.Vendor,
			Latitude:       service.Latitude,
			Longitude:      service.Longitude,
		}

		searchResult = append(searchResult, res)
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": searchResult,
	})
}

func (h *serviceHandler) HandleGetDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	// Get service  info
	service, err := h.server.DB.FindServiceById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if tags, err := h.server.DB.FindAllTagByServiceId(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		service.Tags = tags
	}

	// Get vendor info
	vendor, err := h.server.DB.FindVendorByService(service.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get count per review rating
	reviews, err := h.server.DB.FindAllReviewByService(service.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	images, err := h.server.DB.FindAllPhotosByServiceId(service.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"serviceInfo":    service,
		"vendorInfo":     vendor,
		"serviceImages":  images,
		"countPerRating": countPerRating,
	})
}

func (h *serviceHandler) HandleGetByVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	services, err := h.server.DB.FindServiceByVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

// takes origin as QueryString ex: origin=lat,long
func (h *serviceHandler) HandleFindRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service, err := h.server.DB.FindServiceById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find service")
	}

	origin, err := parseOrigin(c.QueryParam("origin"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid origin coordinates")
	}

	distination := models.NewLocationWithData(service.Latitude, service.Longitude)

	polyline, err := h.server.RouteEngine.FindRoute(origin, distination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find routes at the moment")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"polyline": polyline,
	})
}

func parseOrigin(query string) (*models.Location, error) {
	coords := strings.Split(query, ",")
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, err
	}

	origin := models.NewLocationWithData(lat, long)

	return origin, nil
}
