package query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

// Returns the first 10 locations and an error if any
func GetLocations() ([]types.Location, error) {
	query := `
        SELECT
            address, longitude, latitude
        FROM 
            Location
        LIMIT 
            10
    `

	var locations []types.Location
	err := db.Connection.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
