package message_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetAcquaintances(userId int) ([]types.Acquaintance, error) {
	query := `
        SELECT DISTINCT 
            u.id,
            u.name,
            u.imageUrl
        FROM 
            User u 
        JOIN 
            Message m ON u.id = m.sender OR u.id = m.receiver
        WHERE
            u.id <> ?
    `

	acquaintances := make([]types.Acquaintance, 0)
	err := db.Connection.Select(&acquaintances, query, userId)
	if err != nil {
		return nil, err
	}

	return acquaintances, nil
}
