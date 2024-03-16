package types

type Review struct {
	Id       int `db:"id"`
	VendorId int `db:"vendorId"`
	Rating   int `db:"rating"`
}

type ReviewCount struct {
	Rating string `db:"rating"`
	Count  int    `db:"count"`
}
