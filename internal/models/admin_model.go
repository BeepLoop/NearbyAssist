package models

type AdminModel struct {
	Model
	UpdateableModel
	Username           string `db:"username"`
	Email              string `db:"email"`
	Password           string `db:"password"`
	UsernameHash       string `db:"usernameHash"`
	EmailHash          string `db:"emailHash"`
	MustChangePassword bool   `db:"mustChangePassword"`
	Suspended          bool   `db:"suspended"`
	Role               string `db:"role"`
}
