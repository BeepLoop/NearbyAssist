package models

import "database/sql"

type UserStatusFilter string

const (
	USER_STATUS_VERIFIED   UserStatusFilter = "verified"
	USER_STATUS_UNVERIFIED UserStatusFilter = "unverified"
	USER_STATUS_ALL        UserStatusFilter = "all"
)

type UserModel struct {
	Model
	Name      string          `json:"name" db:"name"`
	Email     string          `json:"email" db:"email"`
	EmailHash string          `json:"emailHash" db:"emailHash"`
	ImageUrl  string          `json:"imageUrl" db:"imageUrl"`
	Address   sql.NullString  `json:"address" db:"address"`
	Phone     sql.NullString  `json:"phone" db:"phone"`
	Latitude  sql.NullFloat64 `json:"latitude" db:"latitude"`
	Longitude sql.NullFloat64 `json:"longitude" db:"longitude"`

	Banned     bool `db:"banned"`
	Restricted bool `db:"restricted"`

	Verified   bool   `json:"verified" db:"verified"`
	VerifiedAt string `json:"-" db:"verifiedAt"`

	// Socials
	Socials []string
}
