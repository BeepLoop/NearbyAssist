package store

import "nearbyassist/internal/models"

type IAdminStore interface {
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	Login(data *models.SessionModel) error

	// Check if refreshToken exists, if exists return nil else return error
	DoesRefreshTokenExists(refreshToken string) error

	// Check if refreshToken is blacklisted, if blacklisted return nil else return error
	IsRefreshTokenBlacklisted(refreshToken string) error
}
