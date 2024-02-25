package query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func NewMessage(message types.Message) error {
	query := `
        INSERT INTO
            Message (sender, reciever, content)
        VALUES
            (:sender, :reciever, :content)
    `

	_, err := db.Connection.NamedExec(query, message)
	if err != nil {
		return err
	}

	return nil
}
