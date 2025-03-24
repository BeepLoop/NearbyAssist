package admin_repo

import "nearbyassist/internal/models"

type AdminRepository interface {
	Create(data *models.AdminModel) error
	GetAll() ([]*models.AdminModel, error)
	GetAllAdmin() ([]*models.AdminModel, error)
	GetAllStaff() ([]*models.AdminModel, error)
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	DoesUsernameExists(usernamehash string) (bool, error)
	ShouldChangePassword(id string) (bool, error)
}
