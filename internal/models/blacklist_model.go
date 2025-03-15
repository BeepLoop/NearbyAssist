package models

type BlacklistModel struct {
	Model
	UpdateableModel
	Token string `json:"token" db:"token"`
}
