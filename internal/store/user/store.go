package user

import "nearbyassist/internal/models"

type UserStore interface {
	CreateUser(user *models.UserModel) (string, error)
	FindByEmailHash(emailHash string) (*models.UserModel, error)
}
