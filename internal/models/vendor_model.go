package models

type VendorModel struct {
	VendorId string `json:"vendorId" db:"vendorId"`
	Rating   string `json:"rating" db:"rating"`
	JoinedAt string `json:"joinedAt" db:"joinedAt"`

	// Additional fields for joins
	User      UserModel
	Expertise []ExpertiseModel
}
