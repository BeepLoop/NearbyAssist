package request

type NewServicePayload struct {
	VendorId    string   `json:"vendorId" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Rate        string   `json:"rate" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
}

type UpdateServicePayload struct {
	VendorId    string   `json:"vendorId" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Rate        string   `json:"rate" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
}
