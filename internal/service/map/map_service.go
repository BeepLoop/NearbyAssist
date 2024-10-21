package map_service

import map_repo "nearbyassist/internal/repository/map"

type Service struct {
	store map_repo.MapRepository
}

func NewService(store map_repo.MapRepository) *Service {
	return &Service{store: store}
}
