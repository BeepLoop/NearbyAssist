package request

type NewTransaction struct {
	VendorId  int `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId  int `json:"clientId" db:"clientId" validate:"required"`
	ServiceId int `json:"serviceId" db:"serviceId" validate:"required"`
}
