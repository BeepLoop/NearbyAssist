package request

type QRSignatureInput struct {
	ClientID  string `json:"clientId" validate:"required"`
	VendorID  string `json:"vendorId" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
}

type QRSignatureVerifyInput struct {
	ClientID  string `json:"clientId" validate:"required"`
	VendorID  string `json:"vendorId" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}
