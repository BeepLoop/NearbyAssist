package tag

import "nearbyassist/internal/models"

type TagStore interface {
	Create(data *models.TagModel) error
	FindById(id string) (*models.TagModel, error)
	FindAll() ([]*models.TagModel, error)
}
