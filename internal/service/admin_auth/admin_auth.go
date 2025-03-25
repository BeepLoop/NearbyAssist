package adminauth_service

import (
	"errors"
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/core"
)

type Service struct {
	adminStore admin_repo.AdminRepository
	encrypt    core.Encryption
	hash       core.Hash
}

func NewService(adminStore admin_repo.AdminRepository, encrypt core.Encryption, hash core.Hash) *Service {
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

	if admin.MustChangePassword {
		decryptedPassword, err := s.encrypt.DecryptString(admin.Password)
		if err != nil {
			return nil, err
		}

		if password != decryptedPassword {
			return nil, errors.New("Invalid credentials")
		}
	} else {
		if !core.IsPasswordMatch(admin.Password, password) {
			return nil, errors.New("Invalid credentials")
		}
	}

	return admin, nil
}
