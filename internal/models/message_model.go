package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type MessageModel struct {
	Model
	Sender   string `json:"sender" db:"sender"`
	Receiver string `json:"receiver" db:"receiver"`
	Content  string `json:"content" db:"content"`
}

func NewMessageModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *MessageModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &MessageModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (m *MessageModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Message (id, sender, receiver, content)
        VALUES
            (:id, :sender, :receiver, :content)
    `

	if _, err := m.Conn.NamedExecContext(ctx, query, m); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return m.Id, nil
}

func (m *MessageModel) UserMessages() ([]MessageModel, error) {
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

	messages := make([]MessageModel, 0)
	if err := m.Conn.SelectContext(
		ctx,
		&messages,
		query,
		m.Sender,
		m.Receiver,
		m.Receiver,
		m.Sender,
	); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return messages, nil
}

// func MessageValueMapFactory(queryParam string) (map[string]int, error) {
// 	queries := strings.Split(queryParam, "&")
// 	if len(queries) != 2 {
// 		return nil, errors.New("missing required field")
// 	}
//
// 	queryValues := make(map[string]int)
// 	for _, query := range queries {
// 		pair := strings.Split(query, "=")
// 		value, err := strconv.Atoi(pair[1])
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		queryValues[pair[0]] = value
// 	}
//
// 	if _, ok := queryValues["sender"]; ok == false {
// 		return nil, errors.New("missing required field")
// 	}
//
// 	if _, ok := queryValues["receiver"]; ok == false {
// 		return nil, errors.New("missing required field")
// 	}
//
// 	return queryValues, nil
// }
