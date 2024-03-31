package user_query

import "nearbyassist/internal/db"

func CountUsers() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM 
            User
    `

	users := 0
	err := db.Connection.Get(&users, query)
	if err != nil {
		return 0, err
	}

	return users, nil
}
