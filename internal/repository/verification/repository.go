package verification_repo

import "nearbyassist/internal/models"

type VerificationRepository interface {
	Create(data *models.IdentityVerificationModel) (string, error)

	GetAll(status string) ([]*models.IdentityVerificationModel, error)
	FindById(id string) (*models.IdentityVerificationModel, error)

	AcceptRequest(requestId string) error
	RejectRequest(requestId, reason string) error
}
