package vendor_repo

import "nearbyassist/internal/models"

type VendorRepository interface {
	FindById(id string) (*models.VendorModel, error)
	GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error)
	GetTags(serviceId string) ([]string, error)
}
