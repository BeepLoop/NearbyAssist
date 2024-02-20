package query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func RegisterUser(user types.User) error {
	query := `
	        INSERT INTO
	            User (name, email, imageUrl)
	        VALUES
                (:name, :email, :imageUrl)
	    `

	_, err := db.Connection.NamedExec(query, user)
	if err != nil {
		return err
	}

	return nil
}
