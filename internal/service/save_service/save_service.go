package save_service

import (
	"nearbyassist/internal/models"
	saved_service_repo "nearbyassist/internal/repository/saved_service"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	savedServiceStore saved_service_repo.SavedServiceRepository
	serviceStore      service_repo.ServiceRepository
	vendorStore       vendor_repo.VendorRepository
	jwt               core.Authenticator
	encrypt           core.Encryption
}

func NewService(savedServiceStore saved_service_repo.SavedServiceRepository, serviceStore service_repo.ServiceRepository, vendorStore vendor_repo.VendorRepository, jwt core.Authenticator, encrypt core.Encryption) *Service {
	return &Service{
		savedServiceStore: savedServiceStore,
		serviceStore:      serviceStore,
		vendorStore:       vendorStore,
		jwt:               jwt,
		encrypt:           encrypt,
	}
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

	return s.savedServiceStore.SaveService(data)
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

	return s.savedServiceStore.UnsaveService(data)
}

func (s *Service) GetSavedServices(bearerToken string) ([]*response.DetailedServiceResponse, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	saves, err := s.savedServiceStore.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	savedServices := make([]*response.DetailedServiceResponse, 0)

	for _, saved := range saves {
		service, err := s.serviceStore.FindById(saved.ServiceId)
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

		for _, extra := range service.Extras {
			if cipher, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = cipher
			}

			if cipher, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = cipher
			}
		}

		reviews, err := s.serviceStore.GetReviews(saved.ServiceId)
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

		vendor, err := s.vendorStore.FindById(service.VendorId)
		if err != nil {
			return nil, err
		}

		if decrypted, err := s.encrypt.DecryptString(vendor.Name); err != nil {
			return nil, err
		} else {
			vendor.Name = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(vendor.Email); err != nil {
			return nil, err
		} else {
			vendor.Email = decrypted
		}

		if vendor.Phone.Valid {
			if decrypted, err := s.encrypt.DecryptString(vendor.Phone.String); err != nil {
				return nil, err
			} else {
				vendor.PhoneString = decrypted
			}
		}

		detailedService := &response.DetailedServiceResponse{
			Service:        service,
			Vendor:         vendor,
			CountPerRating: countPerRating,
		}

		savedServices = append(savedServices, detailedService)
	}

	return savedServices, nil
}
