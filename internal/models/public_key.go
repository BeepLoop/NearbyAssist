package models

type PublicKeyModel struct {
	Model
	UpdateableModel
	Owner string `json:"owner" db:"owner"`
	Pem   string `json:"pem" db:"pem"`
}
