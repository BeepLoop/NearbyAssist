package store

import "nearbyassist/internal/models"

type IAdminStore interface {
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	Login(data *models.SessionModel) error
}
