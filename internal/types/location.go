package types

type Location struct {
	Address string `query:"address"`
	Point   string `query:"point"`
}

type LocationRegister struct {
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}
