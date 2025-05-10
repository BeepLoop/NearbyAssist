package user_repo

import "nearbyassist/internal/models"

type UserRepository interface {
	CreateUser(user *models.UserModel) (string, error)
	ChangeAddress(userId string, address *models.AddressModel) error
	FindById(id string) (*models.UserModel, error)
	FindByEmailHash(emailHash string) (*models.UserModel, error)

	GetAddress(userId string) (*models.AddressModel, error)
	GetIdentification(userId string) (*models.IdentificationModel, error)

	GetAllUserAccounts(limit, offset int) ([]*models.UserModel, error)
	GetBasicUserAccounts(limit, offset int) ([]*models.UserModel, error)

	IsVendor(userId string) (bool, error)

	GetExpertise(userId string) ([]*models.ExpertiseModel, error)

	AddSocial(data *models.SocialModel) (string, error)
	DeleteSocial(userId, id string) error
	GetSocials(userId string) ([]*models.SocialModel, error)

	IsBanned(userId string) (bool, error)
	BanUser(userId string) error
	UnbanUser(userId string) error

	// isRestricted, isExpired, error
	IsRestricted(userId string) (bool, bool, error)
	RestrictUser(data *models.RestrictionModel) error
	LiftRestrictionIfExpired(userId string) error
	ForceLiftRestriction(userId string) error
}
