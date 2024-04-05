package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetUser(userId int) (*types.User, error) {
	query := `
        SELECT
            id, name, email, imageUrl
        FROM 
            User 
        WHERE
            id = ?
    `

	userResult := new(types.User)
	err := db.Connection.Get(userResult, query, userId)
	if err != nil {
		return nil, err
	}

	return userResult, nil
}
