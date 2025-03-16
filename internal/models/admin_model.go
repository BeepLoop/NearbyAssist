package models

type AdminModel struct {
	Model
	UpdateableModel
	Username     string `json:"username" db:"username"`
	Password     string `json:"-" db:"password"`
	UsernameHash string `json:"-" db:"usernameHash"`
}
