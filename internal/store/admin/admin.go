package store

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type AdminStore struct {
	db *sqlx.DB
}

func NewAdminStore(db *sqlx.DB) *AdminStore {
	return &AdminStore{db: db}
}

func (s *AdminStore) Create(data *models.AdminModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Admin (id, username, password, usernameHash, role) VALUES (:id, :username, :password, :usernameHash, :role)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *AdminStore) Login(data *models.SessionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Session (id, refreshToken) VALUES (:id, :refreshToken)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *AdminStore) FindById(id string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	admin := new(models.AdminModel)

	query := "SELECT id, username, password, role FROM Admin WHERE id = ?"
	if err := s.db.GetContext(ctx, admin, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (s *AdminStore) FindByUsernameHash(hash string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	admin := new(models.AdminModel)

	query := "SELECT id, username, password, role FROM Admin WHERE usernameHash = ?"
	if err := s.db.GetContext(ctx, admin, query, hash); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}
