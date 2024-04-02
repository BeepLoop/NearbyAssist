package vendor_query

import "nearbyassist/internal/db"

func RestrictVendor(vendorId int) error {
	query := `
        UPDATE
            Vendor 
        SET 
            restricted = 1
        WHERE 
            vendorId = ?
    `

	_, err := db.Connection.Exec(query, vendorId)
	if err != nil {
		return err
	}

	return nil
}
