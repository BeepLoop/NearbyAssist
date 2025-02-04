package application_repo

import "nearbyassist/internal/models"

type ApplicationRepository interface {
	Create(data *models.ApplicationModel) (string, error)
	FindById(id string) (*models.ApplicationModel, error)
	NewProof(data *models.ApplicationProofModel) (string, error)
	NewPoliceClearance(data *models.PoliceClearanceModel) (string, error)

	GetAll(status string) ([]*models.ApplicationModel, error)

	AcceptRequest(applicationId string) error
	RejectRequest(applicationId string) error
}
