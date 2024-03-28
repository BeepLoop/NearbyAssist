package transaction_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func ClientOngoingTransactions(clientId int) ([]types.Transaction, error) {
	query := `
        SELECT
            id, vendorId, clientId, serviceId, status
        FROM 
            Transaction 
        WHERE
            status = 'ongoing' AND clientId = ?
    `

	transactions := make([]types.Transaction, 0)
	err := db.Connection.Select(&transactions, query, clientId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
