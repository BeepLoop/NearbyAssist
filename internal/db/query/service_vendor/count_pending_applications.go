package vendor_query

import "nearbyassist/internal/db"

func CountPendingApplications() (int, error) {
    query := `
        SELECT 
            COUNT(*)
        FROM
            Application
        WHERE
            status = 'pending'
    `

    applications := 0
    err := db.Connection.Get(&applications, query)
    if err != nil {
        return 0, err
    }

    return applications, nil
}
