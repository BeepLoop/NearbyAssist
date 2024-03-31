package transaction_query

import "nearbyassist/internal/db"

func CancelledTransactions() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM 
            Transaction 
        WHERE
            status = 'cancelled'
    `

	cancelled := 0
	err := db.Connection.Get(&cancelled, query)
	if err != nil {
		return 0, err
	}

	return cancelled, nil
}
