package vendor_query

import "nearbyassist/internal/db"

func UnrestrictVendor(vendorId int) error {
	query := `
        UPDATE
            Vendor 
        SET 
            restricted = 0
        WHERE 
            vendorId = ?
    `

	_, err := db.Connection.Exec(query, vendorId)
	if err != nil {
		return err
	}

	return nil
}
