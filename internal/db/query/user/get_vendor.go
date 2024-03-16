package user_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetVendor(vendorId string) (*types.VendorData, error) {
	query := `
        SELECT
            v.vendorId, v.rating, u.name, role
        FROM
            Vendor v
        LEFT JOIN User u on u.id = v.vendorId
        WHERE
            vendorId = ?
    `

	vendor := new(types.VendorData)
	err := db.Connection.Get(vendor, query, vendorId)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}
