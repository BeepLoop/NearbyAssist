package service_service

import (
	"cmp"
	"fmt"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	store   service_repo.ServiceRepository
	encrypt auth.Encryption
	hash    auth.Hash
	jwt     auth.Authenticator
	suggest suggestion_engine.Engine
	route   route_engine.Engine
}

func NewService(store service_repo.ServiceRepository, encrypt auth.Encryption, hash auth.Hash, jwt auth.Authenticator, suggest suggestion_engine.Engine, route route_engine.Engine) *Service {
	return &Service{store: store, encrypt: encrypt, hash: hash, jwt: jwt, suggest: suggest, route: route}
}

func (s *Service) CreateService(req *request.NewServicePayload) (string, error) {
	if err := s.store.IsVendor(req.VendorId); err != nil {
		return "", err
	}

	encryptedTitle, err := s.encrypt.EncryptString(req.Title)
	if err != nil {
		return "", err
	}

	encryptedDesc, err := s.encrypt.EncryptString(req.Description)
	if err != nil {
		return "", err
	}

	// Compute signature
	toSign := fmt.Sprintf("%s_%s_%f_%f", req.VendorId, req.Description, req.Latitude, req.Longitude)
	signature, err := s.hash.Generate([]byte(toSign))
	if err != nil {
		return "", err
	}

	if service, err := s.store.FindBySignature(signature); err == nil && service != nil {
		return "", err
	}

	newService := new(models.ServiceModel)
	newService.VendorId = req.VendorId
	newService.Title = encryptedTitle
	newService.Description = encryptedDesc
	newService.Rate = req.Rate
	newService.Tags = req.Tags
	newService.Latitude = req.Latitude
	newService.Longitude = req.Longitude
	newService.Signature = signature

	serviceId, err := s.store.Create(newService)
	if err != nil {
		return "", err
	}

	return serviceId, nil
}

func (s *Service) GetService(serviceId string) (map[string]interface{}, error) {
	service, err := s.store.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	if cipher, err := s.encrypt.DecryptString(service.Title); err != nil {
		return nil, err
	} else {
		service.Title = cipher
	}

	if cipher, err := s.encrypt.DecryptString(service.Description); err != nil {
		return nil, err
	} else {
		service.Description = cipher
	}

	extras := make([]models.ExtraModel, 0)
	for _, extra := range service.Extras {
		decryptedTitle, err := s.encrypt.DecryptString(extra.Title)
		if err != nil {
			return nil, err
		}

		decryptedDescription, err := s.encrypt.DecryptString(extra.Description)
		if err != nil {
			return nil, err
		}

		extras = append(extras, models.ExtraModel{
			Model:           extra.Model,
			UpdateableModel: extra.UpdateableModel,
			Title:           decryptedTitle,
			Description:     decryptedDescription,
			Price:           extra.Price,
		})
	}
	service.Extras = extras

	if tags, err := s.store.GetTags(serviceId); err != nil {
		return nil, err
	} else {
		service.Tags = tags
	}

	reviews, err := s.store.GetReviews(serviceId)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	vendor, err := s.store.GetVendorInfo(service.VendorId)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(vendor.Vendor); err != nil {
		return nil, err
	} else {
		vendor.Vendor = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(vendor.Email); err != nil {
		return nil, err
	} else {
		vendor.Email = decrypted
	}

	vendorData := struct {
		Id           string   `json:"id"`
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		ImageUrl     string   `json:"imageUrl"`
		Rating       string   `json:"rating"`
		IsRestricted int      `json:"isRestricted"`
		Expertise    []string `json:"expertise"`
	}{
		Id:           vendor.VendorId,
		Name:         vendor.Vendor,
		Email:        vendor.Email,
		ImageUrl:     vendor.ImageUrl,
		Rating:       vendor.Rating,
		IsRestricted: vendor.Restricted,
		Expertise:    vendor.Expertise,
	}

	data := map[string]interface{}{
		"serviceInfo":    service,
		"vendorInfo":     vendorData,
		"serviceImages":  photos,
		"countPerRating": countPerRating,
	}

	return data, nil
}

func (s *Service) UpdateService(bearerToken, serviceId string, req *request.UpdateServicePayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if vendor, err := s.store.GetVendorInfo(userId); err != nil {
		return err
	} else {
		if vendor.VendorId != userId {
			return err
		}
	}

	updatedService := new(models.ServiceModel)
	updatedService.Id = req.Id
	updatedService.VendorId = req.VendorId
	updatedService.Rate = req.Rate
	updatedService.Tags = req.Tags
	updatedService.Latitude = req.Latitude
	updatedService.Longitude = req.Longitude

	if cipher, err := s.encrypt.EncryptString(req.Title); err != nil {
		return err
	} else {
		updatedService.Title = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Description); err != nil {
		return err
	} else {
		updatedService.Description = cipher
	}

	// Recompute signature
	toSign := fmt.Sprintf("%s_%s_%s_%f_%f", req.VendorId, req.Title, req.Description, req.Latitude, req.Longitude)
	if signature, err := s.hash.Generate([]byte(toSign)); err != nil {
		return err
	} else {
		updatedService.Signature = signature
	}

	if err := s.store.Update(updatedService); err != nil {
		return err
	}

	return nil
}

func (s *Service) SearchService(params map[string]string) ([]*response.SearchResult, error) {
	services, err := s.store.GeoSpatialSearch(params)
	if err != nil {
		return nil, err
	}

	// Compute service distance
	for _, service := range services {
		origin := new(models.GeoSpatialModel)
		if location, ok := params["l"]; ok {
			if err := origin.FromString(location); err != nil {
				return nil, err
			}
		}

		destination := new(models.GeoSpatialModel)
		destination.Latitude = service.Latitude
		destination.Longitude = service.Longitude

		if distance, err := s.route.GetDistance(origin, destination); err != nil {
			return nil, err
		} else {
			service.Distance = distance
		}
	}

	// Compute service suggestibility score
	scoredServices, err := s.suggest.GenerateSuggestions(services)
	if err != nil {
		return nil, err
	}

	// Rank services
	slices.SortFunc(scoredServices, func(a, b *response.SearchResult) int {
		return cmp.Compare(b.Score, a.Score)
	})
	for i, service := range scoredServices {
		service.Rank = i + 1
	}

	return scoredServices, nil
}

func (s *Service) FindRoute(serviceId string, origin string) (route_engine.PolylineCode, error) {
	service, err := s.store.FindById(serviceId)
	if err != nil {
		return "", err
	}

	lat, long, err := utils.ParseCoordinate(origin)
	if err != nil {
		return "", err
	}

	from := &models.GeoSpatialModel{
		Latitude:  lat,
		Longitude: long,
	}
	distination := &models.GeoSpatialModel{
		Latitude:  service.Latitude,
		Longitude: service.Longitude,
	}

	polyline, err := s.route.GetPolyline(from, distination)
	if err != nil {
		return "", err
	}

	return polyline, nil
}
