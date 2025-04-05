package request

type NewBookingPayload struct {
	VendorId  string  `json:"vendorId" validate:"required"`
	ClientId  string  `json:"clientId" validate:"required"`
	ServiceId string  `json:"serviceId" validate:"required"`
	Cost      string  `json:"cost" validate:"required"`
	Extras    []Extra `json:"extras" validate:"required"`
}
