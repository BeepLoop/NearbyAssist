package session_query

import (
	"nearbyassist/internal/db"
)

func NewSession(username string, token string) error {

	tx, err := db.Connection.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE Session SET status = 'offline' WHERE userId = (SELECT id FROM User where name = ?) AND status = 'online'", username)
	if err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return rollbackError
		}

		return err
	}

	_, err = tx.Exec("INSERT INTO Session (userId, token) VALUES ((SELECT id FROM User WHERE name = ?), ?)", username, token)
	if err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return rollbackError
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
