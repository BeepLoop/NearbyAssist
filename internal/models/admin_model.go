package models

type AdminRole string

const (
	ADMIN_ROLE_ADMIN AdminRole = "admin"
	ADMIN_ROLE_STAFF AdminRole = "staff"
)

type AdminModel struct {
	Model
	UpdateableModel
	Username     string    `json:"username" db:"username" validate:"required"`
	Password     string    `json:"password" db:"password" validate:"required"`
	Role         AdminRole `json:"role" db:"role"`
	UsernameHash string    `json:"usernameHash" db:"usernameHash"`
}

func NewAdminModel() *AdminModel {
	return &AdminModel{}
}
