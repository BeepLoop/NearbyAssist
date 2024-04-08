package models

import (
	"context"
	"errors"
	"nearbyassist/internal/db"
	"strconv"
	"strings"
	"time"
)

type MessageModel struct {
	Model
	UpdateableModel
	Sender   int    `json:"sender" db:"sender"`
	Receiver int    `json:"receiver" db:"receiver"`
	Content  string `json:"content" db:"content"`
}

func NewMessageModel() *MessageModel {
	return &MessageModel{}
}

func MessageModelFactory(queryParam string) (*MessageModel, error) {
	queries := strings.Split(queryParam, "&")
	if len(queries) != 2 {
		return nil, errors.New("missing required field")
	}

	queryValues := make(map[string]int)
	for _, query := range queries {
		pair := strings.Split(query, "=")
		value, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, err
		}

		queryValues[pair[0]] = value
	}

	if _, ok := queryValues["sender"]; ok == false {
		return nil, errors.New("missing required field")
	}

	if _, ok := queryValues["receiver"]; ok == false {
		return nil, errors.New("missing required field")
	}

	return &MessageModel{
		Sender:   queryValues["sender"],
		Receiver: queryValues["receiver"],
	}, nil
}

func (m *MessageModel) Create() (int, error) {
	return 0, nil
}

func (m *MessageModel) Update(id int) error {
	return nil
}

func (m *MessageModel) Delete(id int) error {
	return nil
}
func (m *MessageModel) GetMessages() ([]MessageModel, error) {
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
	err := db.Connection.SelectContext(
		ctx,
		&messages,
		query,
		m.Sender,
		m.Receiver,
		m.Receiver,
		m.Sender,
	)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return messages, nil
}

func (m *MessageModel) GetConversations(userId int) ([]UserModel, error) {
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

	acquaintances := make([]UserModel, 0)
	err := db.Connection.SelectContext(ctx, &acquaintances, query, userId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return acquaintances, nil
}
