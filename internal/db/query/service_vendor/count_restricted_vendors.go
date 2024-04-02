package vendor_query

import "nearbyassist/internal/db"

func CountRestrictedVendors() (int, error) {
	query := `
        SELECT
            COUNT(*)
        FROM
            Vendor 
        WHERE
            restricted = 1
    `

	restricted := 0
	err := db.Connection.Get(&restricted, query)
	if err != nil {
		return 0, err
	}

	return restricted, nil
}
