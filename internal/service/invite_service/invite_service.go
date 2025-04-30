package invite_service

import (
	"context"
	"database/sql"
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

const (
	ERR_DUPLICATE_USERNAME = "username already taken"
	ERR_DUPLICATE_EMAIL    = "email already registered"
	ERR_EXPIRED_INVITE     = "invitation code already expired"
	ERR_CODE_NOT_FOUND     = "invite code does not found"
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
	usernameHash := utils.Must(s.hash.Generate([]byte(username)))
	if doesExists, err := s.adminStore.DoesUsernameExists(usernameHash); err != nil {
		return err
	} else {
		if doesExists {
			return errors.New(ERR_DUPLICATE_USERNAME)
		}
	}

	emailHash := utils.Must(s.hash.Generate([]byte(email)))
	if doesExists, err := s.adminStore.DoesEmailExists(emailHash); err != nil {
		return err
	} else {
		if doesExists {
			return errors.New(ERR_DUPLICATE_EMAIL)
		}
	}

	// NOTE: Password is encrypted with AES256 instead of BCrypt in order to
	// decrypt it and send an email to the user
	encryptedDefaultPassword := utils.Must(s.encrypt.EncryptString(defaultPassword))

	d, err := utils.StringDaysToDuration(duration)
	if err != nil {
		return err
	}
	expiryDate := time.Now().Add(d)
	invitation := &models.InvitationModel{
		Username:     utils.Must(s.encrypt.EncryptString(username)),
		UsernameHash: usernameHash,
		Email:        utils.Must(s.encrypt.EncryptString(email)),
		EmailHash:    emailHash,
		Code:         utils.GenerateId(),
		Password:     encryptedDefaultPassword,
		ExpiredAt:    utils.FormatDateTime(expiryDate),
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
		if err == sql.ErrNoRows {
			return errors.New(ERR_CODE_NOT_FOUND)
		}

		return err
	}

	if isExpired, err := s.inviteStore.IsExpired(invitation.Id); err != nil {
		return err
	} else {
		if isExpired {
			return errors.New(ERR_EXPIRED_INVITE)
		}
	}

	if err := s.inviteStore.Accept(invitation.Id); err != nil {
		return err
	}

	payload := &mailer.InviteAcceptedPayload{
		Username:  utils.Must(s.encrypt.DecryptString(invitation.Username)),
		Password:  utils.Must(s.encrypt.DecryptString(invitation.Password)),
		Email:     utils.Must(s.encrypt.DecryptString(invitation.Email)),
		LoginLink: fmt.Sprintf("%s/auth/login", config.Instance.DOMAIN),
	}

	if err := s.mailer.Send(context.Background(), payload); err != nil {
		return err
	}

	return nil
}
