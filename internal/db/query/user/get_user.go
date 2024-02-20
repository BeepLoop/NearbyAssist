package query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetUser(user types.User) (*types.User, error) {
	query := `
        SELECT
            name, email
        FROM 
            User 
        WHERE
            name = ?
        AND
            email = ?
    `

	userResult := new(types.User)
	err := db.Connection.Get(userResult, query, user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	return userResult, nil
}
