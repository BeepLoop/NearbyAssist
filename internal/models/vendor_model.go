package models

type VendorStatusFilter string

const (
	VENDOR_STATUS_RESTRICTED   VendorStatusFilter = "restricted"
	VENDOR_STATUS_UNRESTRICTED VendorStatusFilter = "unrestricted"
	VENDOR_STATUS_ALL          VendorStatusFilter = "all"
)

type VendorModel struct {
	Model
	UpdateableModel
	VendorId   string `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	Job        string `json:"job" db:"job"`
	Restricted int    `json:"restricted" db:"restricted"`

	// Additional fields for joins
	Vendor   string `json:"vendor" db:"vendor"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
}
