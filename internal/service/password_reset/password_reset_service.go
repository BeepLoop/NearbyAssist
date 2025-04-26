package passwordreset_service

import (
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
