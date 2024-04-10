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

func NewUserModel(db *db.DB) *UserModel {
	return &UserModel{
		Model: Model{Db: db},
	}
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

	res, err := u.Db.Conn.NamedExecContext(ctx, query, u)
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
	err := u.Db.Conn.GetContext(ctx, user, query, id)
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
	err := u.Db.Conn.SelectContext(ctx, &users, query)
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
	err := u.Db.Conn.GetContext(ctx, user, query, email)
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
	err := u.Db.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
