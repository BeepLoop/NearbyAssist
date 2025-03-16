package admin_repo

import "nearbyassist/internal/models"

type AdminRepository interface {
	Create(data *models.AdminModel) error
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	ShouldChangePassword(id string) (bool, error)
	RequestPasswordReset(data *models.PasswordResetRequestModel) error
}
