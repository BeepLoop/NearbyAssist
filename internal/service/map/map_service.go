package map_service

import (
	"nearbyassist/internal/models"
	map_repo "nearbyassist/internal/repository/map"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/service/core"
)

type Service struct {
	mapStore     map_repo.MapRepository
	serviceStore service_repo.ServiceRepository
	encrypt      core.Encryption
}

func NewService(mapStore map_repo.MapRepository, serviceStore service_repo.ServiceRepository, encrypt core.Encryption) *Service {
	return &Service{
		mapStore:     mapStore,
		serviceStore: serviceStore,
		encrypt:      encrypt,
	}
}

func (s *Service) GetServices(query string) ([]*models.ServiceModel, error) {
	services, err := s.serviceStore.FindAllByTag(query)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if decrypted, err := s.encrypt.DecryptString(service.Title); err != nil {
			return nil, err
		} else {
			service.Title = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(service.Description); err != nil {
			return nil, err
		} else {
			service.Description = decrypted
		}
	}

	return services, nil
}
