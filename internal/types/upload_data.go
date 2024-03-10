package types

type UploadData struct {
	VendorId  int    `json:"vendorId" db:"vendorId" validate:"required"`
	ServiceId int    `json:"serviceId" db:"serviceId" validate:"required"`
	ImageUrl  string `db:"imageUrl"`
}
