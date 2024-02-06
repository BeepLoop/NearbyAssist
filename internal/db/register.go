package db

import "nearbyassist/internal/types"

func RegisterUser(user types.User) error {

	query := `
        INSERT INTO
            User (name, email)
        VALUES
            (?, ?)
    `

	_, err := DB_CONN.Exec(query, user.Name, user.Email)
	if err != nil {
		return err
	}

	return nil
}
