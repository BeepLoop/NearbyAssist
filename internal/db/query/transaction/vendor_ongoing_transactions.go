package transaction_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func VendorOngoingTransactions(vendorId int) ([]types.Transaction, error) {
	query := `
        SELECT
            id, vendorId, clientId, serviceId, status
        FROM 
            Transaction 
        WHERE
            status = 'ongoing' AND vendorId= ?
    `

	transactions := make([]types.Transaction, 0)
	err := db.Connection.Select(&transactions, query, vendorId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
