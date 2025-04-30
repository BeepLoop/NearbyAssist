package admin_repo

import "nearbyassist/internal/models"

type AdminRepository interface {
	Create(data *models.AdminModel) error
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	GetAll() ([]*models.AdminModel, error)
	GetAllWithRole(role string) ([]*models.AdminModel, error)
	DoesUsernameExists(usernamehash string) (bool, error)
	DoesEmailExists(emailHash string) (bool, error)
	Suspend(id string) error
	Unsuspend(id string) error
}
