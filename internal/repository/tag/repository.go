package tag_repo

import "nearbyassist/internal/models"

type TagRepository interface {
	Create(data *models.TagModel) error
	FindById(id string) (*models.TagModel, error)
	FindAll() ([]*models.TagModel, error)

	FindAllWithExpertise() ([]*models.ExpertiseModel, error)
}
