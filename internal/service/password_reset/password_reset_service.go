package passwordreset_service

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	passwordreset_repo "nearbyassist/internal/repository/password_reset"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/mailer"
)

type Service struct {
	adminStore         admin_repo.AdminRepository
	passwordResetStore passwordreset_repo.PasswordResetRepository
	mailer             mailer.Mailer
	encrypt            core.Encryption
	hash               core.Hash
}

func NewService(adminStore admin_repo.AdminRepository, passwordResetStore passwordreset_repo.PasswordResetRepository, mailer mailer.Mailer, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		adminStore:         adminStore,
		passwordResetStore: passwordResetStore,
		mailer:             mailer,
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

func (s *Service) ChangePassword(username, oldPassword, newPassword, confirmNewPassword string) error {
	if newPassword != confirmNewPassword {
		return errors.New("password mismatch")
	}

	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return err
	}

	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	// NOTE: This should be guaranteed that admin password is not encrypted with
	// BCrypt because the only way to get here is on first-time login via invite
	// or password reset by admin
	decryptedPassword, err := s.encrypt.DecryptString(admin.Password)
	if err != nil {
		return err
	}

	if oldPassword != decryptedPassword {
		return errors.New("invalid password")
	}

	if !core.IsPasswordSecure(newPassword) {
		return errors.New("insecure password")
	}

	// Encrypt with BCrypt
	encryptedPassword, err := core.BcryptPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.passwordResetStore.ChangePassword(admin.Id, encryptedPassword); err != nil {
		return err
	}

	// TODO: Send email notification to user

	return nil
}

func (s *Service) FulfillResetPassword(requestId, newPassword, confirmationUsername, confirmationPassword string) error {
	request, err := s.passwordResetStore.FindById(requestId)
	if err != nil {
		return err
	}

	usernameHash, err := s.hash.Generate([]byte(confirmationUsername))
	if err != nil {
		return err
	}

	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	if !core.IsPasswordMatch(admin.Password, confirmationPassword) {
		return errors.New("Invalid credentials")
	}

	if !core.IsPasswordSecure(newPassword) {
		return errors.New("insecure password")
	}

	// NOTE: Encrypt the password with AES256 instead of BCrypt in order to
	// decrypt it later for email sending purposes
	encryptedPassword, err := s.encrypt.EncryptString(newPassword)
	if err != nil {
		return err
	}

	if err := s.passwordResetStore.ResetPassword(requestId, encryptedPassword); err != nil {
		return err
	}

	// TODO: Send email to user for their reset credentials
	account, err := s.adminStore.FindById(request.AdminId)
	if err != nil {
		return err
	}

	username, err := s.encrypt.DecryptString(account.Username)
	if err != nil {
		return err
	}

	email, err := s.encrypt.DecryptString(account.Email)
	if err != nil {
		return err
	}

	payload := &mailer.PasswordResetPayload{
		Username:  username,
		Password:  newPassword,
		Email:     email,
		LoginLink: fmt.Sprintf("%s/auth/login", config.Instance.DOMAIN),
	}

	return s.mailer.Send(context.Background(), payload)
}

func (s *Service) RejectResetPassword(requestId, reason string) error {
	// NOTE: Implement reason
	return s.passwordResetStore.Delete(requestId)
}
