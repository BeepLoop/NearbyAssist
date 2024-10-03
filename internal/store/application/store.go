package application

import "nearbyassist/internal/models"

type ApplicationStore interface {
	Create(data *models.ApplicationModel) (string, error)
	FindApplication(id string) (*models.ApplicationModel, error)
	NewProof(data *models.ApplicationProofModel) (string, error)
	NewPoliceClearance(data *models.PoliceClearanceModel) (string, error)
}
