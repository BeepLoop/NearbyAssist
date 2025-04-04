package dashboard_repo

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

type DashboardRepository interface {
	GetUserData() (*models.UserData, error)
	GetBugReportData() (*models.WeeklyBugReportData, error)
	GetVendorReportData() (*models.WeeklyVendorReportData, error)
	GetIdentityVerificationRequestsData() (*models.IdentityVerificationRequestData, error)
	GetVendorApplicationRequestsData() (*models.VendorApplicationRequestData, error)
	GetBookingData() (*models.WeeklyBookingData, error)
}
