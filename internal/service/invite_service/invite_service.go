package invite_service

import (
	"errors"
	"time"

	"nearbyassist/internal/models"
	invitation_repo "nearbyassist/internal/repository/invitation"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	inviteStore invitation_repo.Repository
	encrypt     core.Encryption
	hash        core.Hash
}

func NewService(inviteStore invitation_repo.Repository, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		inviteStore: inviteStore,
		encrypt:     encrypt,
		hash:        hash,
	}
}

func (s *Service) Invite(username, email, duration string) error {
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

	invitation := &models.InvitationModel{
		Username:     encryptedUsername,
		UsernameHash: usernameHash,
		Email:        encryptedEmail,
		EmailHash:    emailHash,
		Code:         utils.GenerateId(),
		ExpiredAt:    utils.FormatDateTime(expiryDate),
	}

	s.inviteStore.Create(invitation)

	// TODO: sent email to invited user

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

	return nil
}
