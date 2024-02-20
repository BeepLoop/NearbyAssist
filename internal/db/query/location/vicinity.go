package query

import (
	"fmt"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

// retrieves all establishments nearby the given position
func SearchVicinity(pos types.Position) ([]types.Location, error) {
	query := fmt.Sprintf(`
        SELECT
            address, ST_AsText(location) as point
        FROM 
            Location 
        WHERE
            ST_Distance_Sphere(
                location,
                ST_GeomFromText('POINT(%s %s)', 4326)
            ) * 0.001 < 10;
    `, pos.Latitude, pos.Longitude)
	// ST_Distance_Sphere returns distance in meters
	// used 0.001 to convert meters to km

	var locations []types.Location
	err := db.Connection.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
