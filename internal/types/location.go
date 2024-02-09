package types

type Location struct {
	Address   string  `query:"address"`
	Longitude float64 `query:"longitude"`
	Latitude  float64 `query:"latitude"`
}
