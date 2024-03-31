package transaction_query

import "nearbyassist/internal/db"

func CountOngoingTransactions() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM
            Transaction 
        WHERE
            status = 'ongoing';
    `

	ongoing := 0
	err := db.Connection.Get(&ongoing, query)
	if err != nil {
		return 0, err
	}

	return ongoing, nil
}
