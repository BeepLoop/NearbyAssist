package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/cache"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/save_service"
	service_service "nearbyassist/internal/service/service"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type serviceHandler struct {
	service_service *service_service.Service
	save_service    *save_service.Service
	resourceService *resource_service.Service
}

func NewHandler(service *service_service.Service, save_service *save_service.Service, resourceService *resource_service.Service) *serviceHandler {
	return &serviceHandler{
		service_service: service,
		save_service:    save_service,
		resourceService: resourceService,
	}
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

	serviceId, err := h.service_service.CreateService(req)
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
	params := c.QueryParams()
	requestURI := c.Request().RequestURI

	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	var serviceDetail *response.DetailedServiceResponse

	if params.Has("fresh") && params.Get("fresh") == "true" {
		detail, err := h.service_service.NewGetService(serviceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error while retrieving service information",
				Error:   err.Error(),
			})
		}

		serviceDetail = detail
	} else {
		inCache, exists := cache.NewGoCache().Get(requestURI)
		if exists {
			serviceDetail = inCache.(*response.DetailedServiceResponse)
		} else {
			detail, err := h.service_service.NewGetService(serviceId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error while retrieving service information",
					Error:   err.Error(),
				})
			}

			serviceDetail = detail
		}
	}

	for _, image := range serviceDetail.Service.Images {
		signedURL, err := h.resourceService.SignURLWithDefaultDuration(image.Url)
		if err != nil {
			c.Logger().Warnf("Error generating service image url: %s\n", err.Error())
			continue
		}
		image.Url = signedURL
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"detail": serviceDetail,
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

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.service_service.UpdateService(bearerToken, serviceId, req); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "you are not allowed to perform this action",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error updating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) AddImage(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID missing",
			Error:   "Service ID is required for this action",
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	imageData, err := h.service_service.AddImage(bearerToken, serviceId, files)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "you are not allowed to do this action",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "could not upload image",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"imageId": imageData.Id,
		"url":     imageData.Url,
	})
}

func (h *serviceHandler) DeleteImage(c echo.Context) error {
	imageId := c.Param("imageId")
	if imageId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Image ID missing",
			Error:   "Image ID is required",
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.service_service.DeleteImage(bearerToken, imageId); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You are not allowed to delete this image",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "could not delete this image",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) AddExtra(c echo.Context) error {
	req := new(request.AddExtraPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "error validating request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	extraId, err := h.service_service.AddExtra(bearerToken, req)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "you are not authorized to do this action",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "could not publish extra",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"extraId": extraId,
	})
}

func (h *serviceHandler) EditExtra(c echo.Context) error {
	req := new(request.EditExtraPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "error validating request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.service_service.EditExtra(bearerToken, req); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "you are not authorized to do this action",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "could not edit service extra",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) DeleteExtra(c echo.Context) error {
	extraId := c.Param("extraId")
	if extraId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Extra ID is missing",
			Error:   "Extra ID is missing",
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.service_service.DeleteExtra(bearerToken, extraId); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "you are not authorized to do this action",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "could not delete service extra",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) SaveService(c echo.Context) error {
	req := new(request.SaveServicePayload)
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

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.save_service.SaveService(bearerToken, req.ServiceId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error saving service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) UnsaveService(c echo.Context) error {
	req := new(request.UnsaveServicePayload)
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

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.save_service.UnsaveService(bearerToken, req.ServiceId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error removing service from saves",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) GetSavedServices(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	services, err := h.save_service.GetSavedServices(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred while getting saved services",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"saves": services,
	})
}

func (h *serviceHandler) SearchService(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())

	services, err := h.service_service.SearchService(params)
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

	polyline, err := h.service_service.FindRoute(serviceId, origin)
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
