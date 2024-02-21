package types

type Location struct {
	OwnerId int    `query:"ownerId" db:"ownerId"`
	Address string `query:"address" db:"address"`
	Point   string `query:"point" db:"location"`
}

type LocationRegister struct {
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}
