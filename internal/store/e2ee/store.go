package e2ee

import "nearbyassist/internal/models"

type E2EEStore interface {
	NewPublicPem(data *models.PublicKeyModel) (string, error)
	GetPublicPem(owner string) (*models.PublicKeyModel, error)

	NewPrivatePem(data *models.PrivateKeyModel) (string, error)
	GetPrivatePem(owner string) (*models.PrivateKeyModel, error)
}
