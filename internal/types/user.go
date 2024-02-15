package types

type User struct {
	Name     string `json:"name" db:"name" validate:"required"`
	Email    string `json:"email" db:"email" validate:"required,email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl" validate:"required"`
}
