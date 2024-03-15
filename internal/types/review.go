package types

type Review struct {
	Id       int `db:"id"`
	VendorId int `db:"vendorId"`
	Rating   int `db:"rating"`
}
