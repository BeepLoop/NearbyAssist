package session_query

import (
	"nearbyassist/internal/db"
)

func NewSession(username string, token string) error {
	query := `
        INSERT INTO
            Session (userId, token)
        VALUES
            ((SELECT id FROM User WHERE name = ?), ?)
    `

	_, err := db.Connection.Exec(query, username, token)
	if err != nil {
		return err
	}

	return nil
}
