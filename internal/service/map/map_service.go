package map_service

import (
	"nearbyassist/internal/models"
	map_repo "nearbyassist/internal/repository/map"
)

type Service struct {
	store map_repo.MapRepository
}

func NewService(store map_repo.MapRepository) *Service {
	return &Service{store: store}
}

func (s *Service) GetServices(query string) ([]*models.ServiceModel, error) {
	return s.store.GetAllByTag(query)
}
