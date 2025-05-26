package request

type ChangeAddressPayload struct {
	Address  string   `json:"address" validate:"required"`
	Location Location `json:"location" validate:"required"`
}

type UpdatePhonePayload struct {
	Phone string `json:"phone" validate:"required"`
}
