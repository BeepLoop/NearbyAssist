package invite_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	invitation_repo "nearbyassist/internal/repository/invitation"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/mailer"
	"nearbyassist/internal/utils"
)

type Service struct {
	adminStore  admin_repo.AdminRepository
	inviteStore invitation_repo.Repository
	mailer      mailer.Mailer
	encrypt     core.Encryption
	hash        core.Hash
}

func NewService(adminStore admin_repo.AdminRepository, inviteStore invitation_repo.Repository, mailer mailer.Mailer, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		adminStore:  adminStore,
		inviteStore: inviteStore,
		mailer:      mailer,
		encrypt:     encrypt,
		hash:        hash,
	}
}

func (s *Service) Invite(username, email, defaultPassword, duration string) error {
	d, err := utils.ParseStringDuration(duration)
	if err != nil {
		return err
	}
	expiryDate := time.Now().Add(d)

	encryptedUsername, err := s.encrypt.EncryptString(username)
	if err != nil {
		return err
	}

	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return err
	}

	encryptedEmail, err := s.encrypt.EncryptString(email)
	if err != nil {
		return err
	}

	emailHash, err := s.hash.Generate([]byte(email))
	if err != nil {
		return err
	}

	// NOTE: Password is encrypted with AES256 instead of BCrypt in order to
	// decrypt it and send an email to the user
	encryptedDefaultPassword, err := s.encrypt.EncryptString(defaultPassword)
	if err != nil {
		return err
	}

	invitation := &models.InvitationModel{
		Username:     encryptedUsername,
		UsernameHash: usernameHash,
		Email:        encryptedEmail,
		EmailHash:    emailHash,
		Code:         utils.GenerateId(),
		Password:     encryptedDefaultPassword,
		ExpiredAt:    utils.FormatDateTime(expiryDate),
	}

	doesExists, err := s.adminStore.DoesUsernameExists(usernameHash)
	if err != nil {
		return err
	}
	if doesExists {
		return errors.New("duplicate username")
	}

	if _, err := s.inviteStore.Create(invitation); err != nil {
		return err
	}

	joinUrl := fmt.Sprintf("%s/admin/invites/join", config.Instance.DOMAIN)
	payload := &mailer.InvitationPayload{
		Username:   username,
		Email:      email,
		JoinURL:    joinUrl,
		InviteCode: invitation.Code,
	}

	if err := s.mailer.Send(context.Background(), payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) Join(code string) error {
	invitation, err := s.inviteStore.FindByCode(code)
	if err != nil {
		return err
	}

	isExpired, err := s.inviteStore.IsExpired(invitation.Id)
	if err != nil {
		return err
	}

	if isExpired {
		return errors.New("invitation expired")
	}

	if err := s.inviteStore.Accept(invitation.Id); err != nil {
		return err
	}

	username, err := s.encrypt.DecryptString(invitation.Username)
	if err != nil {
		return err
	}

	password, err := s.encrypt.DecryptString(invitation.Password)
	if err != nil {
		return err
	}

	email, err := s.encrypt.DecryptString(invitation.Email)
	if err != nil {
		return err
	}

	payload := &mailer.InviteAcceptedPayload{
		Username:  username,
		Password:  password,
		Email:     email,
		LoginLink: fmt.Sprintf("%s/auth/login", config.Instance.DOMAIN),
	}

	if err := s.mailer.Send(context.Background(), payload); err != nil {
		return err
	}

	return nil
}
