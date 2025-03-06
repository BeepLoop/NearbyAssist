package models

type VendorComplaintModel struct {
	Model
	UpdateableModel
	VendorId string `json:"vendorId" db:"vendorId"`
	Title    string `json:"title" db:"title"`
	Content  string `json:"content" db:"content"`
}
