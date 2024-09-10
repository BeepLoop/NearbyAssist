package user

import "nearbyassist/internal/models"

type UserStore interface {
	CreateUser(user *models.UserModel) (string, error)
	FindById(id string) (*models.UserModel, error)
	FindByEmailHash(emailHash string) (*models.UserModel, error)
	Login(data *models.SessionModel) error

	// Check if refreshToken exists, if exists return nil else return error
	DoesRefreshTokenExists(refreshToken string) error

	// Check if refreshToken is blacklisted, if blacklisted return nil else return error
	IsRefreshTokenBlacklisted(refreshToken string) error
}
