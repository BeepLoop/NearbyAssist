package models

type UserModel struct {
	Model
	UpdateableModel
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
}

func NewUserModel() *UserModel {
	return &UserModel{}
}
