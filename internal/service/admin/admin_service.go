package admin_service

import (
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/core"
)

type Service struct {
	adminStore admin_repo.AdminRepository
	encrypt    core.Encryption
}

func NewService(adminStore admin_repo.AdminRepository, encrypt core.Encryption) *Service {
	return &Service{
		adminStore: adminStore,
		encrypt:    encrypt,
	}
}

func (s *Service) GetAll() ([]*models.AdminModel, error) {
	accounts, err := s.adminStore.GetAll()
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if plain, err := s.encrypt.DecryptString(account.Username); err != nil {
			return nil, err
		} else {
			account.Username = plain
		}

		if plain, err := s.encrypt.DecryptString(account.Email); err != nil {
			return nil, err
		} else {
			account.Email = plain
		}
	}

	return accounts, nil
}
