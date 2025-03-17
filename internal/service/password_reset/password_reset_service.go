package passwordreset_service

import (
	"errors"
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	passwordreset_repo "nearbyassist/internal/repository/password_reset"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	adminStore         admin_repo.AdminRepository
	passwordResetStore passwordreset_repo.PasswordResetRepository
	encrypt            auth.Encryption
	hash               auth.Hash
}

func NewService(adminStore admin_repo.AdminRepository, passwordResetStore passwordreset_repo.PasswordResetRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{
		adminStore:         adminStore,
		passwordResetStore: passwordResetStore,
		encrypt:            encrypt,
		hash:               hash,
	}
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

	if _, err := s.passwordResetStore.Create(data); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetResetRequest(requestId string) (*models.PasswordResetRequestModel, error) {
	return s.passwordResetStore.FindById(requestId)
}

func (s *Service) GetResetRequests() ([]*models.PasswordResetRequestModel, error) {
	requests, err := s.passwordResetStore.GetAll()
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		if decrypted, err := s.encrypt.DecryptString(request.Username); err != nil {
			return nil, err
		} else {
			request.Username = decrypted
		}
	}

	return requests, nil
}

func (s *Service) ResetPassword(requestId, newPassword, confirmationUsername, confirmationPassword string) error {
	usernameHash, err := s.hash.Generate([]byte(confirmationUsername))
	if err != nil {
		return err
	}

	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	if !auth.IsPasswordMatch(admin.Password, confirmationPassword) {
		return errors.New("Invalid credentials")
	}

	if !auth.IsPasswordSecure(newPassword) {
		return errors.New("insecure password")
	}

	encryptedPassword, err := auth.BcryptPassword(newPassword)
	if err != nil {
		return err
	}

	return s.passwordResetStore.ResetPassword(requestId, encryptedPassword)
}

func (s *Service) RejectResetPassword(requestId, reason string) error {
	// NOTE: Implement reason
	return s.passwordResetStore.Delete(requestId)
}
