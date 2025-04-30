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
	Name      string `db:"name"`
	Email     string `db:"email"`
	EmailHash string `db:"emailHash"`
	ImageUrl  string `db:"imageUrl"`
	Phone     string `db:"phone"`
	Verified  bool   `db:"verified"`

	Banned     bool `db:"banned"`
	Restricted bool `db:"restricted"`

	VerifiedAt sql.NullString `db:"verifiedAt"`

	// Socials
	Socials        []SocialModel
	Address        AddressModel
	Identification IdentificationModel
}
