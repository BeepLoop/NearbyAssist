package vendor_repo

import "nearbyassist/internal/models"

type VendorRepository interface {
	GetAll(limit, offset int) ([]*models.VendorModel, error)
	FindById(id string) (*models.VendorModel, error)
	GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error)
	GetTags(serviceId string) ([]*models.TagModel, error)

	IsRestricted(userId string) (bool, error)
}
