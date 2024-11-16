package verification_repo

import "nearbyassist/internal/models"

type VerificationRepository interface {
	Create(data *models.IdentityVerificationModel) (string, error)

	GetAll() ([]*models.IdentityVerificationModel, error)
}
