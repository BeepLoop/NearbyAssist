package vendor_query

import "nearbyassist/internal/db"

func CountApprovedApplications() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM
            Application 
        WHERE
            status = 'approved'
    `

	count := 0
	err := db.Connection.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}
