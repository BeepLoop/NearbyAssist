package verification

import "nearbyassist/internal/models"

type VerificationStore interface {
	Create(data *models.IdentityVerificationModel) (string, error)
}
