package tag_service

import (
	"nearbyassist/internal/models"
	tag_repo "nearbyassist/internal/repository/tag"
)

type Service struct {
	store tag_repo.TagRepository
}

func NewService(store tag_repo.TagRepository) *Service {
	return &Service{store: store}
}

func (s *Service) GetTags() ([]*models.TagModel, error) {
	tags, err := s.store.FindAll()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *Service) GetExpertise() ([]*models.ExpertiseModel, error) {
	return s.store.FindAllWithExpertise()
}
