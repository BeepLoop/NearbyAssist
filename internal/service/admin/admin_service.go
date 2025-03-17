package admin_service

import (
	"errors"

	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	adminStore admin_repo.AdminRepository
	encrypt    auth.Encryption
	hash       auth.Hash
}

func NewService(adminStore admin_repo.AdminRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{
		adminStore: adminStore,
		encrypt:    encrypt,
		hash:       hash,
	}
}

func (s *Service) Login(username, password string) (*models.AdminModel, error) {
	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return nil, err
	}

	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return nil, err
	}

	if auth.IsPasswordMatch(admin.Password, password) == false {
		return nil, errors.New("Invalid credentials")
	}

	return admin, nil
}
