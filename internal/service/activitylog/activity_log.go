package activitylog

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	instance *activityLog
)

const (
	ACTION_LOGIN                    = "LOGIN"
	ACTION_LOGOUT                   = "LOGOUT"
	ACTION_APPROVED_VERIFICATION    = "APPROVED_VERIFICATION"
	ACTION_REJECTED_VERIFICATION    = "REJECTED_VERIFICATION"
	ACTION_APPROVED_APPLICATION     = "APPROVED_APPLICATION"
	ACTION_REJECTED_APPLICATION     = "REJECTED_APPLICATION"
	ACTION_CREATED_EXPERTISE        = "CREATED_EXPERTISE"
	ACTION_ADDED_EXPERTISE_TAG      = "ADDED_EXPERTISE_TAG"
	ACTION_VIEWED_USER              = "VIEWED_USER"
	ACTION_VIEWED_VENDOR            = "VIEWED_VENDOR"
	ACTION_BANNED_USER              = "BANNED_USER"
	ACTION_UNBANNED_USER            = "UNBANNED_USER"
	ACTION_SUSPENDED_USER           = "SUSPENDED_USER"
	ACTION_UNSUSPENDED_USER         = "UNSUSPENDED_USER"
	ACTION_VIEWED_REPORT            = "VIEWED_REPORT"
	ACTION_RESOLVED_REPORT          = "RESOLVED_REPORT"
	ACTION_DISMISSED_REPORT         = "DISMISSED_REPORT"
	ACTION_SENT_INVITE              = "SENT_INVITE"
	ACTION_PASSWORD_RESET_REQUEST   = "PASSWORD_RESET_REQUEST"
	ACTION_FULFILLED_PASSWORD_RESET = "FULFILLED_PASSWORD_RESET"
	ACTION_REJECTED_PASSWORD_RESET  = "REJECTED_PASSWORD_RESET"
	ACTION_PASSWORD_CHANGE          = "PASSWORD_CHANGE"
	ACTION_RESETTED_SSE             = "RESETTED_SSE"
	ACTION_NOTIFIED_VENDORS         = "NOTIFIED_VENDORS"
)

type Input struct {
	AdminId    string `db:"adminId"`
	Action     string `db:"action"`
	TargetType string `db:"targetType"` // "user" | "admin" | "other"
	TargetId   string `db:"targetId"`
}

type activityLog struct {
	db      *sqlx.DB
	encrypt core.Encryption
}

func NewActivityLogService(db *sqlx.DB, encrypt core.Encryption) *activityLog {
	if instance != nil {
		return instance
	}

	instance = &activityLog{
		db:      db,
		encrypt: encrypt,
	}
	return instance
}

func MustGetInstance() *activityLog {
	if instance == nil {
		panic("activity log instance not initialized")
	}

	return instance
}

func (a *activityLog) Create(input Input) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ActivityLog(id, adminId, action, targetType)
        VALUES
            (?, ?, ?, 'none')
    `

	input.Action = utils.Must(a.encrypt.EncryptString(strings.ToUpper(input.Action)))
	if _, err := a.db.ExecContext(ctx, query, utils.GenerateId(), input.AdminId, input.Action); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (a *activityLog) CreateWithTarget(input Input) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ActivityLog(id, adminId, action, targetType, targetId)
        VALUES
            (?, ?, ?, ?, ?)
    `

	input.Action = utils.Must(a.encrypt.EncryptString(strings.ToUpper(input.Action)))
	if _, err := a.db.ExecContext(ctx, query, utils.GenerateId(), input.AdminId, input.Action, input.TargetType, input.TargetId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (a *activityLog) GetAll(limit, offset int, query string) ([]*models.ActivityLogModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getLogs := "SELECT * FROM ActivityLog ORDER BY createdAt DESC LIMIT ? OFFSET ?"

	logs := make([]*models.ActivityLogModel, 0)
	if err := a.db.SelectContext(ctx, &logs, getLogs, limit, offset); err != nil {
		return nil, err
	}

	adminUsername := "SELECT username FROM Admin WHERE id = ?"
	targetUserEmail := "SELECT email FROM User WHERE id = ?"
	for _, log := range logs {
		var username string
		if err := a.db.GetContext(ctx, &username, adminUsername, log.AdminId); err != nil {
			return nil, err
		}
		log.AdminUsername = utils.Must(a.encrypt.DecryptString(username))

		if log.TargetType == "admin" {
			if log.TargetId.Valid {
				var username string
				if err := a.db.GetContext(ctx, &username, adminUsername, log.TargetId.String); err != nil {
					return nil, err
				}
				log.Target = utils.Must(a.encrypt.DecryptString(username))
			}
		} else if log.TargetType == "user" {
			if log.TargetId.Valid {
				var email string
				if err := a.db.GetContext(ctx, &email, targetUserEmail, log.TargetId.String); err != nil {
					return nil, err
				}
				log.Target = utils.Must(a.encrypt.DecryptString(email))
			}
		}

		log.Action = utils.Must(a.encrypt.DecryptString(log.Action))
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return logs, nil
}
