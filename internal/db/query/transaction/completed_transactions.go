package transaction_query

import "nearbyassist/internal/db"

func CompletedTransactions() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM 
            Transaction 
        WHERE
            status = 'done'
    `

	completed := 0
	err := db.Connection.Get(&completed, query)
	if err != nil {
		return 0, err
	}

	return completed, nil
}
