package tag_service

import (
	"nearbyassist/internal/models"
	expertise_repo "nearbyassist/internal/repository/expertise"
	tag_repo "nearbyassist/internal/repository/tag"
)

type Service struct {
	tagStore       tag_repo.TagRepository
	expertiseStore expertise_repo.ExpertiseRepository
}

func NewService(tagStore tag_repo.TagRepository, expertiseStore expertise_repo.ExpertiseRepository) *Service {
	return &Service{tagStore: tagStore, expertiseStore: expertiseStore}
}

func (s *Service) GetTags() ([]*models.TagModel, error) {
	tags, err := s.tagStore.FindAll()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *Service) GetExpertiseList() ([]*models.ExpertiseModel, error) {
	return s.expertiseStore.GetAll()
}
