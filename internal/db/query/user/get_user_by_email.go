package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetUserByEmail(email string) (*types.User, error) {
	query := `
        SELECT
            id, name, email, imageUrl
        FROM
            User
        WHERE
            email = ?
    `
	user := new(types.User)
	err := db.Connection.Get(user, query, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
