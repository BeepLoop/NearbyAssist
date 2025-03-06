package admin_repo

import "nearbyassist/internal/models"

type AdminRepository interface {
	Create(data *models.AdminModel) error
	FindById(id string) (*models.AdminModel, error)
	FindByUsernameHash(hash string) (*models.AdminModel, error)
	Login(data *models.SessionModel) error

	// Sets the refreshToken in session to offline and adds the refreshToken to blacklist
	Logout(refreshToken string) error

	// Check if refreshToken exists, if exists return nil else return error
	DoesRefreshTokenExists(refreshToken string) error

	// Check if refreshToken is blacklisted, if blacklisted return nil else return error
	IsRefreshTokenBlacklisted(refreshToken string) error
}
