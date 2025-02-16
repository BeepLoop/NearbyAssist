package management_service

import (
	"nearbyassist/internal/models"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	store   user_repo.UserRepository
	encrypt auth.Encryption
}

func NewService(store user_repo.UserRepository, encrypt auth.Encryption) *Service {
	return &Service{store: store, encrypt: encrypt}
}

func (s *Service) GetUsers(limit, offset int) ([]*models.UserModel, error) {
	accounts, err := s.store.GetAllUserAccounts(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if decrypted, err := s.encrypt.DecryptString(account.Name); err != nil {
			return nil, err
		} else {
			account.Name = decrypted
		}
	}

	return accounts, nil
}
