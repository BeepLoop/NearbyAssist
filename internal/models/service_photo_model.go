package models

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId"`
	VendorId  string `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}
