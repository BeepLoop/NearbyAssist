package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

// Returns the first 10 locations and an error if any
func GetServices() ([]types.Service, error) {
	query := `
        SELECT
            vendor, title, description, rate, ST_AsText(location) as location, category
        FROM 
            Service
        LIMIT
            10
    `

	var locations []types.Service
	err := db.Connection.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
