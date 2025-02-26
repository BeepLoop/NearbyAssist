package models

import "database/sql"

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
	Restricted int    `json:"restricted" db:"restricted"`

	// Additional fields for joins
	Vendor    string         `json:"vendor" db:"vendor"`
	Email     string         `json:"email" db:"email"`
	Phone     sql.NullString `db:"phone"`
	ImageUrl  string         `json:"imageUrl" db:"imageUrl"`
	Expertise []string       `json:"expertise" db:"expertise"`
	Socials   []string
}
