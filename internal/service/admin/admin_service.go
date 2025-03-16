package admin_service

import (
	"errors"

	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	adminStore repository.AdminRepository
	encrypt    auth.Encryption
	hash       auth.Hash
}

func NewService(adminStore repository.AdminRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
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

func (s *Service) RequestPasswordReset(username string) error {
	if username == "" {
		return errors.New("invalid username")
	}

	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return err
	}

	account, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	data := &models.PasswordResetRequestModel{
		AdminId: account.Id,
	}

	return s.adminStore.RequestPasswordReset(data)
}
