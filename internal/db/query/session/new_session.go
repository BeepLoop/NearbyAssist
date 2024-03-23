package session_query

import (
	"nearbyassist/internal/db"
)

func NewSession(username, email, token string) error {
	tx, err := db.Connection.Beginx()
	if err != nil {
		return err
	}

	var onlineCount int
	err = tx.Get(&onlineCount, "SELECT COUNT(*) FROM Session WHERE userId = (SELECT id FROM User WHERE name = ? AND email = ?) AND status = 'online'", username, email)
	if err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return rollbackError
		}

		return err
	}

	if onlineCount > 0 {
		_, err = tx.Exec("UPDATE Session SET status = 'offline' WHERE userId = (SELECT id FROM User where name = ? AND email = ?) AND status = 'online'", username, email)
		if err != nil {
			if rollbackError := tx.Rollback(); rollbackError != nil {
				return rollbackError
			}

			return err
		}
	}

	_, err = tx.Exec("INSERT INTO Session (userId, token) VALUES ((SELECT id FROM User WHERE name = ? AND email = ?), ?)", username, email, token)
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
