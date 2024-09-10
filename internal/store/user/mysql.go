package user

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlUserStore struct {
	db *sqlx.DB
}

func NewMysqlUserStore(db *sqlx.DB) *MysqlUserStore {
	return &MysqlUserStore{
		db: db,
	}
}

func (s *MysqlUserStore) CreateUser(user *models.UserModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := store.GenerateNanoId()
	if err != nil {
		return "", err
	}
	user.Id = id

	query := "INSERT INTO User (id, name, email, imageUrl, emailHash) VALUES (:id, :name, :email, :imageUrl, :hash)"
	if _, err := s.db.NamedExecContext(ctx, query, user); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return id, nil
}

func (s *MysqlUserStore) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)

	query := "SELECT id, name, email, imageUrl, verified FROM User WHERE emailHash = ?"
	if err := s.db.GetContext(ctx, user, query, emailHash); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}
