package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) NewUser(user *models.UserModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO User (name, email, imageUrl) VALUES (:name, :email, :imageUrl)"

	res, err := m.Conn.NamedExecContext(ctx, query, user)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) FindUserById(id int) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl FROM User WHERE id = ?"

	user := models.NewUserModel()
	err := m.Conn.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (m *Mysql) FindUserByEmail(email string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl FROM User WHERE email = ?"

	user := models.NewUserModel()
	err := m.Conn.GetContext(ctx, user, query, email)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (m *Mysql) CountUser() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM User"

	count := -1
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return count, nil
}
