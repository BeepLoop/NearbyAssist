package management

import "nearbyassist/internal/models"

type ManagementStore interface {
	CreateStaff(data *models.AdminModel) (string, error)

	GetUserById(id string) (*models.UserModel, error)

	RestrictVendor(id string) error
	UnrestrictVendor(id string) error

	GetApplications(filter map[string]string) ([]*models.ApplicationModel, error)
	ApproveApplication(id string) error
	RejectApplication(id string) error

	GetTransaction(id string) (*models.TransactionModel, error)

	GetSystemComplaints(filter map[string]string) ([]*models.ComplaintModel, error)
	GetSystemComplaintById(id string) (*models.ComplaintModel, error)

	GetVerificationRequests(filter map[string]string) ([]*models.IdentityVerificationModel, error)
	GetVerificationRequestById(id string) (*models.IdentityVerificationModel, error)
}
