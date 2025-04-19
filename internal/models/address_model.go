package models

type AddressModel struct {
	Id        string  `db:"id"`
	Address   string  `db:"address"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}
