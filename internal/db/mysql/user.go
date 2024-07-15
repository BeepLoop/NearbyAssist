package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) NewUser(user *models.UserModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO User (name, email, imageUrl, emailHash) VALUES (:name, :email, :imageUrl, :hash)"

	res, err := m.Conn.NamedExecContext(ctx, query, user)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) CheckUserVerification(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT verified FROM User WHERE id = ?"

	var verified bool
	if err := m.Conn.GetContext(ctx, &verified, query, id); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return false, nil
}

func (m *Mysql) FindUserById(id int) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl, verified FROM User WHERE id = ?"

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

func (m *Mysql) FindUserByEmailHash(hash string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl FROM User WHERE emailHash = ?"

	user := models.NewUserModel()
	if err := m.Conn.GetContext(ctx, user, query, hash); err != nil {
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

	count := 0
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
