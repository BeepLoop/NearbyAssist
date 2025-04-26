package passwordreset_service

import (
	"errors"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

func (s *Service) ChangePassword(username, oldPassword, newPassword, confirmNewPassword string) error {
	if newPassword != confirmNewPassword {
		return errors.New(ERR_MISMATCHING_PASSWORD)
	}

	usernameHash := utils.Must(s.hash.Generate([]byte(username)))
	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	// NOTE: This should be guaranteed that admin password is not encrypted with
	// BCrypt because the only way to get here is on first-time login via invite
	// or password reset by admin
	decryptedPassword := utils.Must(s.encrypt.DecryptString(admin.Password))
	if oldPassword != decryptedPassword {
		return errors.New(ERR_INVALID_CREDENTIALS)
	}

	if !core.IsPasswordSecure(newPassword) {
		return errors.New(ERR_WEAK_PASSWORD)
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
