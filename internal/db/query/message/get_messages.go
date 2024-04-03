package message_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetMessages(message types.Message) ([]types.Message, error) {
	query := `
        SELECT
            id, sender, receiver, content
        FROM 
            Message
        WHERE
            sender = ? AND receiver = ? 
        OR
            sender = ? AND receiver = ?
        ORDER BY
            createdAt
    `

	var messages []types.Message
	err := db.Connection.Select(
		&messages,
		query,
		message.Sender,
		message.Receiver,
		message.Receiver,
		message.Sender,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
