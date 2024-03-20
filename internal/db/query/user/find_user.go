package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func FindUser(name, email string) (*types.User, error) {
	query := `
        SELECT 
            id, name, email, imageUrl
        FROM
            User 
        WHERE
            name = ? AND email = ?
    `

	user := new(types.User)
	err := db.Connection.Get(user, query, name, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
