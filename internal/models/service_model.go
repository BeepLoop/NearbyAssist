package models

type ServiceModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	VendorId    int    `json:"vendorId" db:"vendorId" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Rate        string `json:"rate" db:"rate" validate:"required"`
}

func NewServiceModel() *ServiceModel {
	return &ServiceModel{}
}
