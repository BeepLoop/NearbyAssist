package save_service

import (
	"nearbyassist/internal/models"
	saved_service_repo "nearbyassist/internal/repository/saved_service"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
)

type Service struct {
	savedStore   saved_service_repo.SavedServiceRepository
	serviceStore service_repo.ServiceRepository
	jwt          auth.Authenticator
	encrypt      auth.Encryption
}

func NewService(store saved_service_repo.SavedServiceRepository, serviceStore service_repo.ServiceRepository, jwt auth.Authenticator, encrypt auth.Encryption) *Service {
	return &Service{savedStore: store, serviceStore: serviceStore, jwt: jwt, encrypt: encrypt}
}

func (s *Service) SaveService(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	data := &models.SavedServiceModel{
		UserId:    userId,
		ServiceId: serviceId,
	}

	return s.savedStore.SaveService(data)
}

func (s *Service) UnsaveService(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	data := &models.SavedServiceModel{
		UserId:    userId,
		ServiceId: serviceId,
	}

	return s.savedStore.UnsaveService(data)
}

func (s *Service) GetSavedServices(bearerToken string) (*response.SavedServicesResponse, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	saves, err := s.savedStore.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	savedServices := &response.SavedServicesResponse{}

	for _, entry := range saves {
		service, err := s.serviceStore.FindById(entry.ServiceId)
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

		if tags, err := s.serviceStore.GetTags(entry.ServiceId); err != nil {
			return nil, err
		} else {
			service.Tags = tags
		}

		reviews, err := s.serviceStore.GetReviews(entry.ServiceId)
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

		photos, err := s.serviceStore.GetPhotos(entry.ServiceId)
		if err != nil {
			return nil, err
		}

		vendor, err := s.serviceStore.GetVendorInfo(service.VendorId)
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
			Id           string `json:"id"`
			Name         string `json:"name"`
			Email        string `json:"email"`
			ImageUrl     string `json:"imageUrl"`
			Rating       string `json:"rating"`
			IsRestricted int    `json:"isRestricted"`
		}{
			Id:           vendor.VendorId,
			Name:         vendor.Vendor,
			Email:        vendor.Email,
			ImageUrl:     vendor.ImageUrl,
			Rating:       vendor.Rating,
			IsRestricted: vendor.Restricted,
		}

		serviceData := response.SavedServiceData{
			Service:        service,
			Vendor:         vendorData,
			Photos:         photos,
			CountPerRating: countPerRating,
		}

		savedServices.Services = append(savedServices.Services, serviceData)
	}

	return savedServices, nil
}
