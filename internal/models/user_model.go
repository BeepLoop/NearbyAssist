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
	UpdateableModel
	Name      string          `json:"name" db:"name"`
	Email     string          `json:"email" db:"email"`
	EmailHash string          `json:"emailHash" db:"emailHash"`
	ImageUrl  string          `json:"imageUrl" db:"imageUrl"`
	Verified  bool            `json:"verified" db:"verified"`
	Address   sql.NullString  `json:"address" db:"address"`
	Latitude  sql.NullFloat64 `json:"latitude" db:"latitude"`
	Longitude sql.NullFloat64 `json:"longitude" db:"longitude"`
}
