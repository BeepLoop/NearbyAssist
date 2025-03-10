package expertise_service

import (
	"nearbyassist/internal/models"
	expertise_repo "nearbyassist/internal/repository/expertise"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	store   expertise_repo.ExpertiseRepository
	encrypt auth.Encryption
	hash    auth.Hash
}

func NewService(store expertise_repo.ExpertiseRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{
		store:   store,
		encrypt: encrypt,
		hash:    hash,
	}
}

func (s *Service) CreateExpertise(data *models.ExpertiseModel) (string, error) {
	return s.store.Create(data)
}

func (s *Service) GetAllExpertise() ([]*models.ExpertiseModel, error) {
	return s.store.GetAll()
}

func (s *Service) FindExpertise(query string) (*models.ExpertiseModel, error) {
	return s.store.FindByTitle(query)
}

func (s *Service) AddTagToExpertise(expertiseId string, data *models.TagModel) (string, error) {
	return s.store.CreateTag(expertiseId, data)
}
