package models

type PasswordResetRequestModel struct {
	Model
	AdminId string `db:"adminId"`

	// Join fields
	Username string `db:"username"`
}
