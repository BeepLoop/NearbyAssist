package types

type Complaint struct {
	Id       int    `db:"id"`
	VendorId int    `db:"vendorId" json:"vendorId" validate:"required"`
	Code     int    `db:"code" json:"code" validate:"required"`
	Title    string `db:"title" json:"title" validate:"required"`
	Content  string `db:"content" json:"content" validate:"required"`
}
