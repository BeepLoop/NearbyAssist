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

func (s *Service) GetAll(filter string) ([]*models.AdminModel, error) {
	accounts := make([]*models.AdminModel, 0)

	if filter == "" || filter == "all" {
		if res, err := s.adminStore.GetAll(); err != nil {
			return nil, err
		} else {
			accounts = res
		}
	} else if filter == "admin" {
		if res, err := s.adminStore.GetAllAdmin(); err != nil {
			return nil, err
		} else {
			accounts = res
		}
	} else if filter == "staff" {
		if res, err := s.adminStore.GetAllStaff(); err != nil {
			return nil, err
		} else {
			accounts = res
		}
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
