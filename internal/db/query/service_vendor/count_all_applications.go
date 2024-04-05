package vendor_query

import "nearbyassist/internal/db"

func CountAllApplications() (int, error) {
	query := `
        SELECT
            COUNT(*)
        FROM
            Application 
    `

	count := 0
	err := db.Connection.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}
