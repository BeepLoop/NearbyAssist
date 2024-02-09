package db

import "nearbyassist/internal/types"

// retrieves all establishments nearby the given position
func SearchVicinity(pos types.Position) ([]types.Location, error) {
	query := `
        SELECT
            * 
        FROM
            addresses
        WHERE
            ST_distance_sphere(
                point(?,?),
                point(longitude, latitude)
            ) * 0.001 < 500
    `
	// ST_distance_sphere returns distance in meters
	// used 0.001 to convert meters to km

	var locations []types.Location
	err := DB_CONN.Select(&locations, query, pos.Longitude, pos.Latitude)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
