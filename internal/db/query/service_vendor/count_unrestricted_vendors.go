package vendor_query

import "nearbyassist/internal/db"

func CountUnrestrictedVendors() (int, error) {
	query := `
        SELECT
            COUNT(*)
        FROM
            Vendor 
        WHERE
            restricted = 0
    `

	unrestricted := 0
	err := db.Connection.Get(&unrestricted, query)
	if err != nil {
		return 0, err
	}

	return unrestricted, nil
}
