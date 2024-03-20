package session_query

import "nearbyassist/internal/db"

func LogoutSession(name, email string) error {
	query := `
        UPDATE 
            Session 
        SET
            status = 'offline'
        WHERE
            userId = (SELECT id FROM User WHERE name = ? AND email = ?)
    `

	_, err := db.Connection.Exec(query, name, email)
	if err != nil {
		return err
	}

	return nil
}
