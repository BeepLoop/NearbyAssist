package admin_service

import (
	"errors"

	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	store   repository.AdminRepository
	encrypt auth.Encryption
	hash    auth.Hash
}

func NewService(store repository.AdminRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{store: store, encrypt: encrypt, hash: hash}
}

func (s *Service) Login(username, password string) (*models.AdminModel, error) {
	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return nil, err
	}

	admin, err := s.store.FindByUsernameHash(usernameHash)
	if err != nil {
		return nil, err
	}

	if auth.IsPasswordMatch(admin.Password, password) == false {
		return nil, errors.New("Invalid credentials")
	}

	return admin, nil
}
