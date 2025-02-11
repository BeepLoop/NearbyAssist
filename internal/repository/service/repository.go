package service_repo

import (
	"nearbyassist/internal/models"
)

type ServiceRepository interface {
	Create(data *models.ServiceModel) (string, error)
	FindAll() ([]*models.ServiceModel, error)
	FindById(id string) (*models.ServiceModel, error)
	FindBySignature(signature string) (*models.ServiceModel, error)

	Update(data *models.ServiceModel) error

	Delete(serviceId string) error

	// return nil if vendorId is found in vendor, else error
	IsVendor(vendorId string) error

	GetVendorInfo(vendorId string) (*models.VendorModel, error)

	GetTags(serviceId string) ([]*models.TagModel, error)

	GetReviews(serviceId string) ([]*models.ReviewModel, error)

	GetPhotos(serviceId string) ([]*models.ServicePhotoModel, error)

	GetAllByVendorId(vendorId string) ([]*models.ServiceModel, error)

	GeoSpatialSearch(params map[string]string) ([]*models.GeoSpatialSearchResult, error)
}
