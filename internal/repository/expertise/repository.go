package expertise_repo

import "nearbyassist/internal/models"

type ExpertiseRepository interface {
	Create(data *models.ExpertiseModel) (string, error)

	GetAll() ([]*models.ExpertiseModel, error)
	FindById(id string) (*models.ExpertiseModel, error)
	FindByTitle(title string) (*models.ExpertiseModel, error)
}
