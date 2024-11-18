package saved_service_repo

import "nearbyassist/internal/models"

type SavedServiceRepository interface {
	SaveService(data *models.SavedServiceModel) error
	UnsaveService(data *models.SavedServiceModel) error
	FindByUserId(userId string) ([]*models.SavedServiceModel, error)
}
