package passwordreset_repo

import "nearbyassist/internal/models"

type PasswordResetRepository interface {
	Create(data *models.PasswordResetRequestModel) (string, error)
	Delete(id string) error
	GetAll() ([]*models.PasswordResetRequestModel, error)
	FindById(id string) (*models.PasswordResetRequestModel, error)
	FindByAdminId(id string) (*models.PasswordResetRequestModel, error)
}
