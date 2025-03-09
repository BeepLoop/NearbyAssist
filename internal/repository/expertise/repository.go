package expertise_repo

import "nearbyassist/internal/models"

type ExpertiseRepository interface {
	Create(data *models.ExpertiseModel) (string, error)
	CreateTag(expertiseId string, data *models.TagModel) (string, error)

	GetAll() ([]*models.ExpertiseModel, error)
	FindById(id string) (*models.ExpertiseModel, error)
}
