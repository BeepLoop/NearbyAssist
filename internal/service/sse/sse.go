package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	instance *SSE
)

type SSE struct {
	verification         int
	application          int
	report               int
	passwordResetRequest int
	bugs                 int
	pendingService       int
}

func New() *SSE {
	if instance != nil {
		return instance
	}

	instance = &SSE{
		verification:         0,
		application:          0,
		report:               0,
		passwordResetRequest: 0,
		bugs:                 0,
		pendingService:       0,
	}
	return instance
}

func (s *SSE) DecreaseVerification() {
	if s.verification <= 0 {
		return
	}

	s.verification--
}

func (s *SSE) IncreaseVerification() {
	s.verification++
}

func (s *SSE) DecreaseApplication() {
	if s.application <= 0 {
		return
	}

	s.application--
}

func (s *SSE) IncreaseApplication() {
	s.application++
}

func (s *SSE) DecreaseReport() {
	if s.report <= 0 {
		return
	}

	s.report--
}

func (s *SSE) IncreaseReport() {
	s.report++
}

func (s *SSE) DecreasePasswordResetRequest() {
	if s.passwordResetRequest <= 0 {
		return
	}

	s.passwordResetRequest--
}

func (s *SSE) IncreasePasswordResetRequest() {
	s.passwordResetRequest++
}

func (s *SSE) DecreaseBugReport() {
	if s.bugs <= 0 {
		return
	}

	s.bugs--
}

func (s *SSE) IncreaseBugReport() {
	s.bugs++
}

func (s *SSE) DecreasePendingService() {
	if s.pendingService <= 0 {
		return
	}

	s.pendingService--
}

func (s *SSE) IncreasePendingService() {
	s.pendingService++
}

func (s *SSE) GetMarshalled() ([]byte, error) {
	data := struct {
		Verification         int `json:"verification"`
		Application          int `json:"application"`
		Report               int `json:"report"`
		PasswordResetRequest int `json:"prr"`
		Bugs                 int `json:"bugs"`
		PendingService       int `json:"pendingService"`
	}{
		Verification:         s.verification,
		Application:          s.application,
		Report:               s.report,
		PasswordResetRequest: s.passwordResetRequest,
		Bugs:                 s.bugs,
		PendingService:       s.pendingService,
	}

	return json.Marshal(data)
}

func (s *SSE) SetValues(db *sqlx.DB) {
	fmt.Println("setting values of sse")

	s.verification = s.queryVerificationCount(db)
	s.application = s.queryApplicationCount(db)
	s.report = s.queryReportedUserCount(db)
	s.passwordResetRequest = s.queryPasswordResetRequestCount(db)
	s.bugs = s.queryBugCount(db)
	s.pendingService = s.queryPendingService(db)

	fmt.Println("setting values complete")
}

func (s *SSE) queryVerificationCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM IdentityVerification WHERE status = 'pending'"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		fmt.Println("error retrieving verification count: ", err.Error())
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
		fmt.Println("error retrieving application count: ", err.Error())
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
		fmt.Println("error retrieving user report count: ", err.Error())
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
		fmt.Println("error retrieving password reset request count: ", err.Error())
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}

func (s *SSE) queryBugCount(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM BugReport WHERE completedAt IS NULL"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		fmt.Println("error retrieving bug report count: ", err.Error())
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}

func (s *SSE) queryPendingService(db *sqlx.DB) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM Service WHERE status = 'under_review'"
	count := 0
	if err := db.GetContext(ctx, &count, query); err != nil {
		fmt.Println("error retrieving pending service count: ", err.Error())
		return 0
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0
	}

	return count
}
