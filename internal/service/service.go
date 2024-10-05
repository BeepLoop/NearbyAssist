package handler

import (
	"cmp"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/store/service"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

type ServiceService struct {
	store     service.ServiceStore
	jwt       auth.Authenticator
	encryptor auth.Encryption
	route     route_engine.Engine
	suggest   suggestion_engine.Engine
}

func NewServiceService(store service.ServiceStore, encryptor auth.Encryption, jwt auth.Authenticator, route route_engine.Engine, suggest suggestion_engine.Engine) *ServiceService {
	return &ServiceService{
		store:     store,
		encryptor: encryptor,
		jwt:       jwt,
		route:     route,
		suggest:   suggest,
	}
}

func (s *ServiceService) GetServices(c echo.Context) error {
	services, err := s.store.FindAll()
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

func (s *ServiceService) Create(c echo.Context) error {
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

	if err := s.store.IsVendor(req.VendorId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "You are not a registered vendor",
			Error:   err.Error(),
		})
	}

	encryptedDesc, err := s.encryptor.EncryptString(req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting description",
			Error:   auth.ENCRYPTION_ERR,
		})
	}

	newService := new(models.ServiceModel)
	newService.VendorId = req.VendorId
	newService.Description = encryptedDesc
	newService.Rate = req.Rate
	newService.Tags = req.Tags
	newService.Latitude = req.Latitude
	newService.Longitude = req.Longitude

	serviceId, err := s.store.Create(newService)
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

func (s *ServiceService) GetService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	service, err := s.store.FindById(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Service not found",
			Error:   err.Error(),
		})
	}

	if cipher, err := s.encryptor.DecryptString(service.Description); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		service.Description = cipher
	}

	if tags, err := s.store.GetTags(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting service tags",
			Error:   err.Error(),
		})
	} else {
		service.Tags = tags
	}

	reviews, err := s.store.GetReviews(serviceId)
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

	photos, err := s.store.GetPhotos(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting service images",
			Error:   err.Error(),
		})
	}

	vendor, err := s.store.GetVendorInfo(service.VendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Vendor not found",
			Error:   err.Error(),
		})
	}

	if decrypted, err := s.encryptor.DecryptString(vendor.Vendor); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		vendor.Vendor = decrypted
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"serviceInfo":    service,
		"vendorInfo":     vendor,
		"serviceImages":  photos,
		"countPerRating": countPerRating,
	})
}

func (s *ServiceService) Update(c echo.Context) error {
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

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	if vendor, err := s.store.GetVendorInfo(userId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting vendor",
			Error:   err.Error(),
		})
	} else {
		if vendor.VendorId != userId {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You do not own this service",
				Error:   "You do not own this service",
			})
		}
	}

	updatedService := new(models.ServiceModel)
	updatedService.Id = req.VendorId
	updatedService.Rate = req.Rate
	updatedService.Tags = req.Tags

	if cipher, err := s.encryptor.EncryptString(req.Description); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting description",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		updatedService.Description = cipher
	}

	if err := s.store.Update(updatedService); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error updating service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ServiceService) Delete(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	if vendor, err := s.store.GetVendorInfo(userId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting vendor",
			Error:   err.Error(),
		})
	} else {
		if vendor.VendorId != userId {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You do not own this service",
				Error:   "You do not own this service",
			})
		}
	}

	if err := s.store.Delete(serviceId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error deleting service",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ServiceService) Search(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())

	services, err := s.store.GeoSpatialSearch(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error searching services",
			Error:   err.Error(),
		})
	}

	// Compute service distance
	for _, service := range services {
		origin := new(models.GeoSpatialModel)
		if location, ok := params["l"]; ok {
			if err := origin.FromString(location); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error Constructing origin",
					Error:   err.Error(),
				})
			}
		}

		destination := new(models.GeoSpatialModel)
		destination.Latitude = service.Latitude
		destination.Longitude = service.Longitude

		if distance, err := s.route.GetDistance(origin, destination); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting distance",
				Error:   err.Error(),
			})
		} else {
			service.Distance = distance
		}
	}

	// Compute service suggestability score
	scoredServices, err := s.suggest.GenerateSuggestions(services)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error generating suggestability score",
			Error:   err.Error(),
		})
	}

	// Rank services
	slices.SortFunc(scoredServices, func(a, b *response.SearchResult) int {
		return cmp.Compare(b.Score, a.Score)
	})
	for i, service := range scoredServices {
		service.Rank = i + 1
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": scoredServices,
	})
}

func (s *ServiceService) GetVendorServices(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	services, err := s.store.GetAllByVendorId(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting services",
			Error:   err.Error(),
		})
	}

	for _, service := range services {
		if plain, err := s.encryptor.DecryptString(service.Description); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: auth.DECRYPTION_ERR,
				Error:   err.Error(),
			})
		} else {
			service.Description = plain
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (s *ServiceService) FindRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID must be a number",
			Error:   "Service ID must be a number",
		})
	}

	service, err := s.store.FindById(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Service not found",
			Error:   err.Error(),
		})
	}

	lat, long, err := utils.ParseCoordinate(c.QueryParam("origin"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid origin",
			Error:   err.Error(),
		})
	}

	origin := &models.GeoSpatialModel{
		Latitude:  lat,
		Longitude: long,
	}
	distination := &models.GeoSpatialModel{
		Latitude:  service.Latitude,
		Longitude: service.Longitude,
	}

	polyline, err := s.route.GetPolyline(origin, distination)
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
