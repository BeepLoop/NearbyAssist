package passwordreset_service

import (
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	passwordreset_repo "nearbyassist/internal/repository/password_reset"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/mailer"
	"nearbyassist/internal/utils"
)

type Service struct {
	adminStore         admin_repo.AdminRepository
	passwordResetStore passwordreset_repo.PasswordResetRepository
	mailer             mailer.Mailer
	encrypt            core.Encryption
	hash               core.Hash
}

func NewService(
	adminStore admin_repo.AdminRepository,
	passwordResetStore passwordreset_repo.PasswordResetRepository,
	mailer mailer.Mailer,
	encrypt core.Encryption,
	hash core.Hash,
) *Service {
	return &Service{
		adminStore:         adminStore,
		passwordResetStore: passwordResetStore,
		mailer:             mailer,
		encrypt:            encrypt,
		hash:               hash,
	}
}

func (s *Service) GetAdminFromUsername(username string) (*models.AdminModel, error) {
	usernameHash := utils.Must(s.hash.Generate([]byte(username)))
	return s.adminStore.FindByUsernameHash(usernameHash)
}

func (s *Service) GetAdminFromRequestId(requestId string) (*models.AdminModel, error) {
	prr, err := s.passwordResetStore.FindById(requestId)
	if err != nil {
		return nil, err
	}

	admin := &models.AdminModel{
		Model: models.Model{Id: prr.AdminId},
	}

	return admin, nil
}
