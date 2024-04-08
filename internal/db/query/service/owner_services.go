package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetOwnerServices(ownerId int) ([]types.Service, error) {
	query := `
        SELECT 
            id, title, description, format(rate, 2) as rate, categoryId as category
        FROM 
            Service 
        WHERE
            vendorId = ?
    `

	services := make([]types.Service, 0)
	err := db.Connection.Select(&services, query, ownerId)
	if err != nil {
		return nil, err
	}

	return services, nil
}
