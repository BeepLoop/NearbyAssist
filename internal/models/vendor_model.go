package models

type VendorStatus string

const (
	VENDOR_STATUS_RESTRICTED   VendorStatus = "restricted"
	VENDOR_STATUS_UNRESTRICTED VendorStatus = "unrestricted"
)

type VendorModel struct {
	Model
	UpdateableModel
	VendorId   int    `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	Job        string `json:"job" db:"job"`
	Restricted int    `json:"restricted" db:"restricted"`
}

func NewVendorModel() *VendorModel {
	return &VendorModel{}
}
