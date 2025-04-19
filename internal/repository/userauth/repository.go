package userauth

import "nearbyassist/internal/models"

type Repository interface {
	Login(data *models.SessionModel) error
	Logout(refreshToken string) error
	FindSessionByToken(refreshToken string) (*models.SessionModel, error)
	IsRefreshTokenBlacklisted(refreshToken string) error
}
