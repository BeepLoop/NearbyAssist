package request

import "nearbyassist/internal/models"

type NewService struct {
	VendorId    int      `json:"vendorId" db:"vendorId" validate:"required"`
	Description string   `json:"description" db:"description" validate:"required"`
	Rate        string   `json:"rate" db:"rate" validate:"required"`
	Tags        []string `json:"tags" db:"tags" validate:"required"`
	models.GeoSpatialModel
}

type UpdateService struct {
	Id          int      `json:"id" db:"id"`
	VendorId    int      `json:"vendorId" db:"vendorId" validate:"required"`
	Description string   `json:"description" db:"description" validate:"required"`
	Rate        string   `json:"rate" db:"rate" validate:"required"`
	Tags        []string `json:"tags" db:"tags" validate:"required"`
	models.GeoSpatialModel
}
