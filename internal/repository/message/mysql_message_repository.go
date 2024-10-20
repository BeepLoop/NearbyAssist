package message_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlMessageRepository struct {
	db *sqlx.DB
}

func NewMysqlChatRepository(db *sqlx.DB) *MysqlMessageRepository {
	return &MysqlMessageRepository{
		db: db,
	}
}
func (s *MysqlMessageRepository) Create(data *models.MessageModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if generatedId, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = generatedId
	}

	query := `
        INSERT INTO
            Message (id, sender, receiver, content)
        VALUES
            (:id, :sender, :receiver, :content)
    `

	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlMessageRepository) GetMessages(user1, user2 string) ([]*models.MessageModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id, sender, receiver, content, createdAt
        FROM
            Message
        WHERE
            sender = ? AND receiver = ?
        OR
            sender = ? AND receiver = ?
        ORDER BY
            createdAt
    `

	messages := make([]*models.MessageModel, 0)
	if err := s.db.SelectContext(ctx, &messages, query, user1, user2, user2, user1); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return messages, nil
}

func (s *MysqlMessageRepository) GetConversations(userId string) ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT DISTINCT
            u.id,
            u.name,
            u.imageUrl
        FROM
            User u
        JOIN
            Message m ON u.id = m.sender OR u.id = m.receiver
        WHERE
            u.id <> ?
    `

	conversations := make([]*models.UserModel, 0)
	if err := s.db.SelectContext(ctx, &conversations, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return conversations, nil
}
