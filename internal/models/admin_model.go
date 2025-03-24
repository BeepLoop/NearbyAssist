package models

type AdminModel struct {
	Model
	UpdateableModel
	Username           string `json:"username" db:"username"`
	Email              string `json:"email" db:"email"`
	Password           string `json:"-" db:"password"`
	UsernameHash       string `json:"-" db:"usernameHash"`
	EmailHash          string `json:"-" db:"emailHash"`
	MustChangePassword bool   `json:"-" db:"mustChangePassword"`
	Role               string `json:"role" db:"role"`
}
