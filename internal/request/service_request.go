package request

import "nearbyassist/internal/models"

type NewService struct {
	models.Model
	VendorId    string   `json:"vendorId" db:"vendorId" validate:"required"`
	Description string   `json:"description" db:"description" validate:"required"`
	Rate        string   `json:"rate" db:"rate" validate:"required"`
	Tags        []string `json:"tags" db:"tags" validate:"required"`
	models.GeoSpatialModel
}

type UpdateService struct {
	Id          string   `json:"id" db:"id"`
	VendorId    string   `json:"vendorId" db:"vendorId" validate:"required"`
	Description string   `json:"description" db:"description" validate:"required"`
	Rate        string   `json:"rate" db:"rate" validate:"required"`
	Tags        []string `json:"tags" db:"tags" validate:"required"`
	models.GeoSpatialModel
}
