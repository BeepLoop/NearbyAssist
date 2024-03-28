package types

type Transaction struct {
	Id        int    `db:"id"`
	VendorId  int    `db:"vendorId"`
	ClientId  int    `db:"clientId"`
	ServiceId int    `db:"serviceId"`
	Status    string `db:"status"`
}
