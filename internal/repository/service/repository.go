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

	// Return nil if vendorId is found in vendor, else error
	IsVendor(vendorId string) error

	GetVendorInfo(vendorId string) (*models.VendorModel, error)

	GetTags(serviceId string) ([]*models.TagModel, error)

	GetReviews(serviceId string) ([]*models.ReviewModel, error)

	FindPhotoById(imageId string) (*models.ServicePhotoModel, error)
	GetPhotos(serviceId string) ([]*models.ServicePhotoModel, error)
	AddImage(data *models.ServicePhotoModel) error
	DeleteImage(imageId string) error

	GetAllByVendorId(vendorId string) ([]*models.ServiceModel, error)

	GeoSpatialSearch(params map[string]string) ([]*models.GeoSpatialSearchResult, error)

	FindExtraById(extraId string) (*models.ExtraModel, error)
	AddExtra(data *models.ExtraModel) error
	EditExtra(data *models.ExtraModel) error
	DeleteExtra(extraId string) error
}
