package types

type ServicePhoto struct {
	ServiceId int    `db:"serviceId"`
	VendorId  int    `db:"vendorId"`
	Url       string `db:"url"`
}
