package expertise_repo

import "nearbyassist/internal/models"

type ExpertiseRepository interface {
	GetAll() ([]*models.ExpertiseModel, error)
	FindById(id string) (*models.ExpertiseModel, error)
}
