package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetServiceDetails(serviceId string) (*types.ServiceDetails, error) {
	query := `
        SELECT 
            s.id,
            vendorId,
            title,
            description,
            u.name,
            u.imageUrl,
            format(s.rate, 2) as rate,
            v.rating,
            v.role as vendorRole
        FROM 
            Service s 
            LEFT JOIN User u ON u.id = s.vendor 
            LEFT JOIN Vendor v ON v.vendorId = s.vendor 
        WHERE 
            s.id = ?
    `

	details := new(types.ServiceDetails)
	err := db.Connection.Get(details, query, serviceId)
	if err != nil {
		return nil, err
	}

	return details, nil
}
