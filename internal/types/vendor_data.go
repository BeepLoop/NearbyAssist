package types

type VendorData struct {
	VendorId int     `db:"vendorId"`
	Name     string  `db:"name"`
	Rating   float64 `db:"rating"`
}
