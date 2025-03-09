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

func (s *Service) GetAllExpertise() ([]*models.ExpertiseModel, error) {
	expertise, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}

	return expertise, nil
}
