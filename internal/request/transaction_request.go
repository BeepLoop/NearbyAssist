package request

import "nearbyassist/internal/models"

type NewTransaction struct {
	models.Model
	VendorId  string `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId  string `json:"clientId" db:"clientId" validate:"required"`
	ServiceId string `json:"serviceId" db:"serviceId" validate:"required"`
	Start     string `json:"start" db:"start" validate:"required"`
	End       string `json:"end" db:"end" validate:"required"`
}
