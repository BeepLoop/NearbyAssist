package sse

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	instance *SSE
)

type SSE struct {
	Verification         int `json:"verification"`
	Application          int `json:"application"`
	Report               int `json:"report"`
	PasswordResetRequest int `json:"prr"`
}

func New() *SSE {
	if instance != nil {
		return instance
	}

	instance = &SSE{}
	return instance
}

func (s *SSE) GetMarshalled() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SSE) SetValues(db *sqlx.DB) {
	s.Verification = s.queryVerificationCount(db)
	s.Application = s.queryApplicationCount(db)
	s.Report = s.queryReportedUserCount(db)
	s.PasswordResetRequest = s.queryPasswordResetRequestCount(db)
}

func (s *SSE) queryVerificationCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM IdentityVerification WHERE status = 'pending'"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}

func (s *SSE) queryApplicationCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM Application WHERE status = 'pending'"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}

func (s *SSE) queryReportedUserCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM UserReport WHERE status = 'pending'"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}

func (s *SSE) queryPasswordResetRequestCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM PasswordResetRequest"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}
