package userauth

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type mysqlUserAuthRepository struct {
	db *sqlx.DB
}

func NewMysqlUserAuthRepository(db *sqlx.DB) *mysqlUserAuthRepository {
	return &mysqlUserAuthRepository{db: db}
}

func (s *mysqlUserAuthRepository) Login(data *models.SessionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	query := "INSERT INTO Session (id, refreshToken) VALUES (:id, :refreshToken)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *mysqlUserAuthRepository) Logout(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateSession := "UPDATE Session SET status = 'offline' WHERE refreshToken = ? AND status = 'online'"
	if _, err := s.db.ExecContext(ctx, updateSession, refreshToken); err != nil {
		return err
	}

	id := utils.GenerateId()
	blacklistToken := `INSERT INTO Blacklist (id, token) VALUES (?, ?)`
	if _, err := tx.ExecContext(ctx, blacklistToken, id, refreshToken); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *mysqlUserAuthRepository) FindSessionByToken(refreshToken string) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	session := new(models.SessionModel)
	query := "SELECT id, status, refreshToken FROM Session WHERE refreshToken = ? AND status = 'online'"
	if err := s.db.GetContext(ctx, session, query, refreshToken); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return session, nil
}

func (s *mysqlUserAuthRepository) IsRefreshTokenBlacklisted(refreshToken string) error {
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
