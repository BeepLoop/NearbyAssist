package verification_repo

import "nearbyassist/internal/models"

type VerificationRepository interface {
	CreateLink(userId string) (string, error)
	// Set newTimestamp = true if not creating new request
	UpdateUserInfo(data *models.IdentityVerificationModel, newTimestamp bool) error
	FindById(id string) (*models.IdentityVerificationModel, error)
	FindByUserId(userId string) (*models.IdentityVerificationModel, error)
	GetAll(status string) ([]*models.IdentityVerificationModel, error)
	AcceptRequest(requestId string) error
	RejectRequest(requestId, reason string) error
}
