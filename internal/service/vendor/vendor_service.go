package vendor_service

import (
	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	store   repository.VendorRepository
	encrypt auth.Encryption
}

func NewService(store repository.VendorRepository, encrypt auth.Encryption) *Service {
	return &Service{store: store, encrypt: encrypt}
}

func (s *Service) GetVendor(vendorId string) (*models.VendorModel, error) {
	vendor, err := s.store.FindById(vendorId)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}

func (s *Service) GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error) {
	services, err := s.store.GetVendorServiceList(vendorId)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if plain, err := s.encrypt.DecryptString(service.Description); err != nil {
			return nil, err
		} else {
			service.Description = plain
		}

		if tags, err := s.store.GetTags(service.Id); err != nil {
			return nil, err
		} else {
			service.Tags = tags
		}
	}

	return services, nil
}
