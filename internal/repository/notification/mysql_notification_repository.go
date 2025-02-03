package notification_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlNotificationRepository struct {
	db *sqlx.DB
}

func NewMysqlNotificationRepository(db *sqlx.DB) *MysqlNotificationRepository {
	return &MysqlNotificationRepository{db: db}
}

func (s *MysqlNotificationRepository) Create(data *models.NotificationModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = id
	}

	insertQuery := `
        INSERT INTO
            Notification (id, recipient, type, title, content)
        VALUES
            (:id, :recipient, :type, :title, :content)
    `
	if _, err := s.db.NamedExecContext(ctx, insertQuery, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlNotificationRepository) FindById(id string) (*models.NotificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := new(models.NotificationModel)
	query := "SELECT * FROM Notification WHERE id = ?"
	if err := s.db.GetContext(ctx, notification, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return notification, nil
}

func (s *MysqlNotificationRepository) UpdateRead(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Notification SET ReadAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlNotificationRepository) GetAllUnreadByRecipient(id string) ([]*models.NotificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notifications := make([]*models.NotificationModel, 0)

	query := `
        SELECT
            *
        FROM
            Notification
        WHERE
            recipient = ? AND readAt IS NULL
        ORDER BY createdAt DESC
    `
	if err := s.db.SelectContext(ctx, &notifications, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return notifications, nil
}
