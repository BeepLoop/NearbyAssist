package user

import (
	"context"
	"errors"
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

func (s *MysqlUserStore) FindById(id string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)
	query := "SELECT id, name, email, imageUrl, verified FROM User WHERE id = ?"
	err := s.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
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

func (s *MysqlUserStore) DoesRefreshTokenExists(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Session WHERE refreshToken = ?"
	if err := s.db.GetContext(ctx, &count, query, refreshToken); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if count > 0 {
		return nil
	}

	return errors.New("refreshToken not found")
}

func (s *MysqlUserStore) IsRefreshTokenBlacklisted(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Blacklist WHERE token = ?"
	if err := s.db.GetContext(ctx, &count, query, refreshToken); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if count > 0 {
		return nil
	}

	return errors.New("refreshToken blacklisted")
}
