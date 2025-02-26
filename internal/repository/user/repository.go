package user_repo

import "nearbyassist/internal/models"

type UserRepository interface {
	CreateUser(user *models.UserModel) (string, error)
	FindById(id string) (*models.UserModel, error)
	FindByEmailHash(emailHash string) (*models.UserModel, error)
	Login(data *models.SessionModel) error

	GetUserAccountPageData(userId string) (*models.UserAccountPageData, error)

	GetAllUserAccounts(limit, offset int) ([]*models.UserModel, error)

	// Sets the refreshToken in session to offline and adds the refreshToken to blacklist
	Logout(refreshToken string) error

	// Check if refreshToken exists, if exists return nil else return error
	FindSessionByToken(refreshToken string) (*models.SessionModel, error)

	// Check if refreshToken is blacklisted, if blacklisted return nil else return error
	IsRefreshTokenBlacklisted(refreshToken string) error

	IsVendor(userId string) (bool, error)

	GetExpertise(userId string) ([]*models.ExpertiseModel, error)

	AddSocial(data *models.SocialModel) error
	DeleteSocial(userId, id string) error
	GetSocials(userId string) ([]*models.SocialModel, error)
}
