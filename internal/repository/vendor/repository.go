package vendor_repo

import "nearbyassist/internal/models"

type VendorRepository interface {
	GetAll(limit, offset int) ([]*models.VendorModel, error)
	FindByEmailHash(emailhash string) (*models.VendorModel, error)
	FindById(id string) (*models.VendorModel, error)
	GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error)
	IsRestricted(userId string) (bool, bool, error)
	AddExpertise(userId, expertiseId, supportingImage string) error
	GetBookingsWithStatus(vendorId, status string) ([]*models.BookingModel, error)
	CompletedBookingCountOfService(vendorId, serviceId string) (int, error)
	HasExpertise(vendorId, expertiseId string) (bool, error)
	GetPoliceClearance(vendorId string) (*models.PoliceClearanceModel, error)
	IsFullyBookedAt(vendorId, schedule string) (bool, error)
	SetDBL(vendorId string, dbl int) error
}
