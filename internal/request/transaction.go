package request

type NewTransactionPayload struct {
	VendorId  string `json:"vendorId" validate:"required"`
	ClientId  string `json:"clientId" validate:"required"`
	ServiceId string `json:"serviceId" validate:"required"`
	Start     string `json:"start" validate:"required"`
	End       string `json:"end" validate:"required"`
}
