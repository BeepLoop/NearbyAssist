package admin_service

import (
	admin_repo "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/service/core"
)

type Service struct {
	adminStore admin_repo.AdminRepository
	encrypt    core.Encryption
}

func NewService(adminStore admin_repo.AdminRepository, encrypt core.Encryption) *Service {
	return &Service{
		adminStore: adminStore,
		encrypt:    encrypt,
	}
}
