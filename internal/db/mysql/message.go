package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) NewMessage(message models.MessageModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Message (sender, receiver, content)
        VALUES
            (:sender, :receiver, :content)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, message)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) GetMessages(senderId, receiverId int) ([]models.MessageModel, error) {
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

	messages := make([]models.MessageModel, 0)
	err := m.Conn.SelectContext(
		ctx,
		&messages,
		query,
		senderId,
		receiverId,
		receiverId,
		senderId,
	)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return messages, nil
}

func (m *Mysql) GetAllUserConversations(userId int) ([]models.UserModel, error) {
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

	conversations := make([]models.UserModel, 0)
	err := m.Conn.SelectContext(ctx, &conversations, query, userId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return conversations, nil
}
