package passwordreset_service

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/config"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/mailer"
	"nearbyassist/internal/utils"
)

type ResetRequestInput struct {
	HandlerAdminId       string
	RequestId            string
	NewPassword          string
	ConfirmationUsername string
	ConfirmationPassword string
}

func (s *Service) FulfillResetPassword(input ResetRequestInput) error {
	request, err := s.passwordResetStore.FindById(input.RequestId)
	if err != nil {
		return err
	}

	// Prevent fulfilling own password reset request
	if request.AdminId == input.HandlerAdminId {
		return errors.New(ERR_SELF_RESETTING_PASSWORD)
	}

	usernameHash := utils.Must(s.hash.Generate([]byte(input.ConfirmationUsername)))
	admin, err := s.adminStore.FindByUsernameHash(usernameHash)
	if err != nil {
		return err
	}

	if !core.IsPasswordMatch(admin.Password, input.ConfirmationPassword) {
		return errors.New(ERR_INVALID_CREDENTIALS)
	}

	if !core.IsPasswordSecure(input.NewPassword) {
		return errors.New(ERR_WEAK_PASSWORD)
	}

	// NOTE: Encrypt the password with AES256 instead of BCrypt in order to
	// decrypt it later for email sending purposes
	encryptedPassword := utils.Must(s.encrypt.EncryptString(input.NewPassword))
	if err := s.passwordResetStore.ResetPassword(input.RequestId, encryptedPassword); err != nil {
		return err
	}

	requestorAccount, err := s.adminStore.FindById(request.AdminId)
	if err != nil {
		return err
	}

	payload := &mailer.PasswordResetPayload{
		Username:  utils.Must(s.encrypt.DecryptString(requestorAccount.Username)),
		Password:  input.NewPassword,
		Email:     utils.Must(s.encrypt.DecryptString(requestorAccount.Email)),
		LoginLink: fmt.Sprintf("%s/auth/login", config.Instance.DOMAIN),
	}

	return s.mailer.Send(context.Background(), payload)
}
