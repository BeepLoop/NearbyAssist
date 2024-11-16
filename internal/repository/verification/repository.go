package verification_repo

import "nearbyassist/internal/models"

type VerificationRepository interface {
	Create(data *models.IdentityVerificationModel) (string, error)

	GetAll() ([]*models.IdentityVerificationModel, error)
	FindById(id string) (*models.IdentityVerificationModel, error)

	AcceptRequest(id string) error
	RejectRequest(id string) error
}
