package admin_repo

import "nearbyassist/internal/models"

type UserStatusFilter string
type VendorStatusFilter string
type ApplicationStatusFilter string

const (
	USER_STATUS_VERIFIED   UserStatusFilter = "verified"
	USER_STATUS_UNVERIFIED UserStatusFilter = "unverified"
	USER_STATUS_ALL        UserStatusFilter = "all"

	VENDOR_STATUS_RESTRICTED   VendorStatusFilter = "restricted"
	VENDOR_STATUS_UNRESTRICTED VendorStatusFilter = "unrestricted"
	VENDOR_STATUS_ALL          VendorStatusFilter = "all"

	APPLICATION_STATUS_PENDING  ApplicationStatusFilter = "pending"
	APPLICATION_STATUS_APPROVED ApplicationStatusFilter = "approved"
	APPLICATION_STATUS_REJECTED ApplicationStatusFilter = "rejected"
)

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

	UserCount(filter UserStatusFilter) (int, error)
	VendorCount(filter VendorStatusFilter) (int, error)
	ApplicationCount(filter ApplicationStatusFilter) (int, error)
	ComplaintCount() (int, error)
}
