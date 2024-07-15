package models

type UserModel struct {
	Model
	UpdateableModel
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
	Verified bool   `json:"verified" db:"verified"`
	Hash     string `json:"hash" db:"hash"`
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func NewUserModelWithData(name, email, imageUrl string) *UserModel {
	return &UserModel{
		Name:     name,
		Email:    email,
		ImageUrl: imageUrl,
	}
}
