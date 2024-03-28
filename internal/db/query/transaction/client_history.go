package transaction_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetClientHistory(userId int) ([]types.TransactionHistory, error) {
	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            s.title as service,
            t.status
        FROM 
            Transaction t
        LEFT JOIN User uVendor ON uVendor.id = t.vendorId
        LEFT JOIN User uClient ON uClient.id = t.clientId
        LEFT JOIN Service s ON s.id = t.serviceId
        WHERE
            status = 'done' OR status = 'cancelled' AND t.clientId = ?
    `

	history := make([]types.TransactionHistory, 0)
	err := db.Connection.Select(&history, query, userId)
	if err != nil {
		return nil, err
	}

	return history, nil
}
