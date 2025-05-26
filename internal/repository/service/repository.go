package service_repo

import (
	"nearbyassist/internal/models"
)

type ServiceRepository interface {
	Create(data *models.ServiceModel) (string, error)
	CreateWithPricingType(data *models.ServiceModel) (string, error)
	FindById(id string) (*models.ServiceModel, error)
	FindBySignature(signature string) (*models.ServiceModel, error)
	GetAll(limit, offset int) ([]*models.ServiceModel, error)
	GetAllUnderReview(limit, offset int) ([]*models.ServiceModel, error)
	GetAllWithTag(tag string) ([]*models.ServiceModel, error)
	GetAllWithTagAny(tags []string) ([]*models.ServiceModel, error)
	FuzzyMatchTags(tags []string) ([]*models.ServiceModel, error)
	GetAllTopRated(limit int) ([]*models.ServiceModel, error)

	Update(data *models.ServiceModel) error
	Resubmit(serviceId string) error

	// Return nil if vendorId is found in vendor, else error
	IsVendor(vendorId string) (bool, error)

	GetReviews(serviceId string) ([]*models.ReviewModel, error)

	FindPhotoById(imageId string) (*models.ServicePhotoModel, error)
	AddImage(data *models.ServicePhotoModel) (string, error)
	DeleteImage(imageId string) error

	IsVendorRestricted(serviceId string) (bool, error)
	IsVendorBanned(serviceId string) (bool, error)

	FindExtraById(extraId string) (*models.ExtraModel, error)
	AddExtra(data *models.ExtraModel) (string, error)
	EditExtra(data *models.ExtraModel) error
	DeleteExtra(extraId string) error

	Disable(serviceId string) error
	Enable(serviceId string) error

	HasActiveBookingWithThisExtra(extraId string) (bool, error)
	HasActiveBookingWithThisService(serviceId string) (bool, error)

	Accept(serviceId string) error
	Reject(serviceId, reason string) error
}
