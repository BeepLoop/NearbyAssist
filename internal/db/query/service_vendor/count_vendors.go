package vendor_query

import "nearbyassist/internal/db"

func CountVendors() (int, error) {
	query := `
        SELECT
            COUNT(*)
        FROM
            Vendor 
    `

	vendors := 0
	err := db.Connection.Get(&vendors, query)
	if err != nil {
		return 0, err
	}

	return vendors, nil
}
