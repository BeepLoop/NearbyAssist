package request

type VerifyAccountPayload struct {
	Name            string  `json:"name" validate:"required"`
	Phone           string  `json:"phone" validate:"required"`
	Address         string  `json:"address" validate:"required"`
	Latitude        float64 `json:"latitude" validate:"required"`
	Longitude       float64 `json:"longitude" validate:"required"`
	IdType          string  `json:"idType" validate:"required"`
	ReferenceNumber string  `json:"referenceNumber" validate:"required"`
}
