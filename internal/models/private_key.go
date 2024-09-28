package models

type PrivateKeyModel struct {
	Model
	UpdateableModel
	Owner string `json:"owner" db:"owner"`
	Pem   string `json:"pem" db:"pem"`
}
