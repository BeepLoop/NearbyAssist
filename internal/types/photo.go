package types

type Photo struct {
	ServiceId int    `db:"serviceId"`
	VendorId  int    `db:"vendorId"`
	Url       string `db:"url"`
}
