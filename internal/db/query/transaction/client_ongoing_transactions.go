package transaction_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func ClientOngoingTransactions(clientId int) ([]types.TransactionData, error) {
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
            status = 'ongoing' AND t.clientId= ?
    `

	transactions := make([]types.TransactionData, 0)
	err := db.Connection.Select(&transactions, query, clientId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
