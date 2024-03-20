package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func RegisterUser(user types.User) (int, error) {
	query := `
	        INSERT INTO
	            User (name, email, imageUrl)
	        VALUES
                (:name, :email, :imageUrl)
	    `

	res, err := db.Connection.NamedExec(query, user)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
