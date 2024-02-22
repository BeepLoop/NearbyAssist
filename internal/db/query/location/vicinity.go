package query

import (
	"fmt"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

// retrieves all establishments nearby the given position
func SearchVicinity(params *types.SearchParams) ([]types.Location, error) {
	query := fmt.Sprintf(`
        SELECT
            ownerId, address, ST_AsText(location) as location
        FROM 
            Location 
        WHERE
            ST_Distance_Sphere(
                location,
                ST_GeomFromText('POINT(%f %f)', 4326)
            ) < ?;
    `, params.Latitude, params.Longitude)
	// ST_Distance_Sphere returns the distance in meters

	var locations []types.Location
	err := db.Connection.Select(&locations, query, params.Radius)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
