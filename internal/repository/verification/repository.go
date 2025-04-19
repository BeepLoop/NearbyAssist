package verification_repo

import "nearbyassist/internal/models"

type VerificationRepository interface {
	Create(userId string) (string, error)
	Update(requestId string, data *models.IdentityVerificationModel) error
	FindById(id string) (*models.IdentityVerificationModel, error)
	FindByUserId(userId string) (*models.IdentityVerificationModel, error)
	GetAll(status string) ([]*models.IdentityVerificationModel, error)
	AcceptRequest(requestId string) error
	RejectRequest(requestId, reason string) error
}
