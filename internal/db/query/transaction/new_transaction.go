package transaction_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func NewTransaction(transaction types.NewTransaction) (int, error) {
	query := `
        INSERT INTO 
            Transaction (vendorId, clientId, serviceId)
        VALUES
            (:vendorId, :clientId, :serviceId)
    `

	res, err := db.Connection.NamedExec(query, transaction)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
