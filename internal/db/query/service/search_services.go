package service_query

import (
	"fmt"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

// retrieves all services matching the query within location radius
func SearchServices(params *types.SearchParams) ([]types.Service, error) {
	query := fmt.Sprintf(`
        SELECT
            id, vendorId, title, description, rate, ST_AsText(location) as location, category
        FROM 
            Service
        WHERE
            title LIKE '%%%s%%'
        AND
            ST_Distance_Sphere(
                location,
                ST_GeomFromText('POINT(%f %f)', 4326)
            ) < ?;
    `, params.Query, params.Latitude, params.Longitude)
	// ST_Distance_Sphere returns the distance in meters

	var locations []types.Service
	err := db.Connection.Select(&locations, query, params.Radius)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
