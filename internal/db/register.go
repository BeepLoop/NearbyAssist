package db

import "nearbyassist/internal/types"

func RegisterUser(user types.User) error {
	query := `
	        INSERT INTO
	            User (name, email)
	        VALUES
            (:name, :email)
	    `

	_, err := DB_CONN.NamedExec(query, user)
	if err != nil {
		return err
	}

	return nil
}
