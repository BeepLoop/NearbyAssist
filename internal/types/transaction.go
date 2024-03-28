package types

type Transaction struct {
	Id        int    `db:"id"`
	VendorId  int    `db:"vendorId"`
	ClientId  int    `db:"clientId"`
	ServiceId int    `db:"serviceId"`
	Status    string `db:"status"`
}

type TransactionData struct {
	Id      int    `db:"id"`
	Vendor  string `db:"vendor"`
	Client  string `db:"client"`
	Service string `db:"service"`
	Status  string `db:"status"`
}
