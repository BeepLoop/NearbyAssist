package request

type NewTransaction struct {
	VendorId  int    `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId  int    `json:"clientId" db:"clientId" validate:"required"`
	ServiceId int    `json:"serviceId" db:"serviceId" validate:"required"`
	Start     string `json:"start" db:"start" validate:"required"`
	End       string `json:"end" db:"end" validate:"required"`
}
