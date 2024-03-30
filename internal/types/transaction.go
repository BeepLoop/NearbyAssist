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

type NewTransaction struct {
	VendorId  int `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId  int `json:"clientId" db:"clientId" validate:"required"`
	ServiceId int `json:"serviceId" db:"serviceId" validate:"required"`
}
