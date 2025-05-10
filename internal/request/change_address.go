package request

type ChangeAddressPayload struct {
	Address  string   `json:"address" validate:"required"`
	Location Location `json:"location" validate:"required"`
}
