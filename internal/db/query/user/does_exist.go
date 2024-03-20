package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func DoesUserExist(user types.User) (bool, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM
            User 
        WHERE 
            name = ?
        AND
            email = ?
    `

	var count int
	err := db.Connection.Get(&count, query, user.Name, user.Email)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
