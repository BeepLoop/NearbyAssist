package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

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

func (u *UserModel) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            User (name, email, imageUrl)
        VALUES
            (:name, :email, :imageUrl)
    `

	res, err := db.Connection.NamedExec(query, u)
	if err != nil {
		return 0, err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(userId), nil
}

func (u *UserModel) Update(id int) error {
	return nil
}

func (u *UserModel) Delete(id int) error {
	return nil
}

func (u *UserModel) FindById(id int) (*UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, name, email, imageUrl
        FROM
            User
        WHERE
            id = ?
    `

	user := new(UserModel)
	err := db.Connection.Get(user, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (u *UserModel) FindAll() ([]UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, name, email, imageUrl 
        FROM
            User 
    `

	users := make([]UserModel, 0)
	err := db.Connection.Select(&users, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return users, nil
}

func (u *UserModel) FindByEmail(email string) (*UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, name, email, imageUrl
        FROM
            User 
        WHERE
            email = ?
    `

	user := new(UserModel)
	err := db.Connection.Get(user, query, email)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (u *UserModel) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            COUNT(*)
        FROM 
            User
    `

	count := 0
	err := db.Connection.Get(&count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
