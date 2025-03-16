package models

type PasswordResetRequestModel struct {
	Model
	AdminId string `db:"adminId"`
}
