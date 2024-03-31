package complaint_query

import "nearbyassist/internal/db"

func CountComplaints() (int, error) {
	query := `
        SELECT 
            COUNT(*)
        FROM
            Complaint
    `

	complaints := 0
	err := db.Connection.Get(&complaints, query)
	if err != nil {
		return 0, err
	}

	return complaints, nil
}
