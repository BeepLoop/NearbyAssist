package query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetMessages(message types.Message) ([]types.Message, error) {
	query := `
        SELECT
            id, sender, reciever, content
        FROM 
            Message
        WHERE
            sender = ? AND reciever = ? 
        OR
            sender = ? AND reciever = ?
        ORDER BY
            createdAt
    `

	var messages []types.Message
	err := db.Connection.Select(
		&messages,
		query,
		message.Sender,
		message.Reciever,
		message.Reciever,
		message.Sender,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
