package models

import "database/sql"

type VendorModel struct {
	VendorId   string `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	JoinedAt   string `json:"joinedAt" db:"joinedAt"`
	Restricted bool   `json:"restricted" db:"restricted"`

	// Additional fields for joins
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Phone     sql.NullString `json:"-" db:"phone"`
	ImageUrl  string         `json:"imageUrl" db:"imageUrl"`
	Expertise []string       `json:"expertise" db:"expertise"`
	Socials   []string       `json:"socials"`
}
