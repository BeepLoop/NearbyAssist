package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetPhotos(serviceId string) ([]types.Photo, error) {
	query := `
        SELECT 
            serviceId, vendorId, url 
        FROM
            Photo 
        WHERE
            serviceId = ?
    `

	photos := make([]types.Photo, 0)
	err := db.Connection.Select(&photos, query, serviceId)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
