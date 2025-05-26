package request

type NewBookingPayload struct {
	VendorId       string  `json:"vendorId" validate:"required"`
	ClientId       string  `json:"clientId" validate:"required"`
	ServiceId      string  `json:"serviceId" validate:"required"`
	RequestedStart string  `json:"requestedStart" validate:"required"`
	RequestedEnd   string  `json:"requestedEnd" validate:"required"`
	Quantity       int     `json:"quantity" validate:"required"`
	Cost           string  `json:"cost" validate:"required"`
	Extras         []Extra `json:"extras" validate:"required"`
}
