package models

type AdminRole string

const (
	ADMIN_ROLE_ADMIN AdminRole = "admin"
	ADMIN_ROLE_STAFF AdminRole = "staff"
)

type AdminModel struct {
	Model
	UpdateableModel
	Username     string    `json:"username" db:"username"`
	Password     string    `json:"-" db:"password"`
	Role         AdminRole `json:"role" db:"role"`
	UsernameHash string    `json:"-" db:"usernameHash"`
}
