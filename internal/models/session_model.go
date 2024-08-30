package models

import (
	"context"
	"errors"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type SessionModel struct {
	Model
	UpdateableModel
	Status       string `json:"status" db:"status"`
	RefreshToken string `json:"refreshToken" db:"refreshToken"`
}

func NewSessionModel(refreshToken string, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *SessionModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &SessionModel{
		Model:        Model{Id: id, Conn: conn},
		RefreshToken: refreshToken,
	}
}

func (s *SessionModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Session (id, token) VALUES (:id, :token)"
	if _, err := s.Conn.NamedExecContext(ctx, query, s); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return s.Id, nil
}

func (s *SessionModel) Logout() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	logout := "UPDATE Session SET status = 'offline' WHERE id = ?"
	if _, err := tx.ExecContext(ctx, logout, s.Id); err != nil {
		return err
	}

	blacklist := "INSERT INTO Blacklist (token) VALUES (?)"
	if _, err := tx.ExecContext(ctx, blacklist, s.RefreshToken); err != nil {
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

func (s *SessionModel) FindByToken() (*SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, token, status FROM Session WHERE token = ?"

	err := s.Conn.GetContext(ctx, s, query, s.RefreshToken)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return s, nil
}

func (s *SessionModel) GetIfActive() (*SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, token, status FROM Session WHERE token = ? AND status = 'online'"

	err := s.Conn.GetContext(ctx, s, query, s.RefreshToken)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return s, nil
}

// func (s *SessionModel) Blacklist() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Blacklist (token) VALUES (?)"
//
// 	_, err := s.Conn.ExecContext(ctx, query, s.RefreshToken)
// 	if err != nil {
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }

func (s *SessionModel) IsBlacklisted() (bool, error) {
	blacklist := NewBlacklistModel(s.IdGenerator, s.Conn)
	if blacklist == nil {
        return true, errors.New(MODEL_INIT_ERROR)
	}

	if result, err := blacklist.FindByToken(s.RefreshToken); err == nil && result != nil {
        return true, nil
	}

    return false, nil
}
