package message_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
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

	data.Id = utils.GenerateId()

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
            createdAt DESC
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

func (s *MysqlMessageRepository) GetConversations(userId string) ([]*models.ConversationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            u.id,
            u.name,
            u.imageUrl,
            m.content AS lastMessage,
            m.createdAt AS lastMessageDate
        FROM 
            User u
            JOIN (
                SELECT 
                    CASE 
                        WHEN sender < receiver THEN sender 
                        ELSE receiver 
                    END AS user1,
                    CASE 
                        WHEN sender < receiver THEN receiver 
                        ELSE sender 
                    END AS user2,
                    MAX(createdAt) AS lastMessageDate
                FROM 
                    Message
                GROUP BY 
                    user1, user2
            ) latest ON (u.id = latest.user1 OR u.id = latest.user2)
            JOIN Message m ON m.createdAt = latest.lastMessageDate
        WHERE 
            u.id <> ?
    `

	conversations := make([]*models.ConversationModel, 0)
	if err := s.db.SelectContext(ctx, &conversations, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return conversations, nil
}
