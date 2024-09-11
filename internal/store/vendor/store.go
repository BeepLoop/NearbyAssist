package vendor

import "nearbyassist/internal/models"

type VendorStore interface {
	FindById(id string) (*models.VendorModel, error)
}
