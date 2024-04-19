package models

type AdminModel struct {
	Model
	UpdateableModel
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
	Role     string `json:"role" db:"role"`
}

func NewAdminModel() *AdminModel {
	return &AdminModel{}
}
