package models

type UserModel struct {
	Model
	UpdateableModel
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
}

func NewUserModel() *UserModel {
	return &UserModel{
		Model:           Model{},
		UpdateableModel: UpdateableModel{},
	}
}

func (u *UserModel) Create() (int, error) {
	return 0, nil
}

func (u *UserModel) Update(id int) error {
	return nil
}

func (u *UserModel) Delete(id int) error {
	return nil
}

func (u *UserModel) FindById(id int) (*UserModel, error) {
	return nil, nil
}

func (u *UserModel) FindAll() ([]UserModel, error) {
	return nil, nil
}

func (u *UserModel) FindByEmail(email string) (*UserModel, error) {
	return nil, nil
}
