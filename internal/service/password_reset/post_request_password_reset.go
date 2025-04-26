package passwordreset_service

import (
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
)

func (s *Service) RequestPasswordReset(username string) error {
	if username == "" {
		return errors.New(ERR_INVALID_USERNAME)
	}

	usernameHash := utils.Must(s.hash.Generate([]byte(username)))
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
