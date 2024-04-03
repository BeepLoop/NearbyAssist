package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetServicePhotos(serviceId int) ([]types.ServicePhoto, error) {
	query := `
        SELECT 
            serviceId, vendorId, url 
        FROM
            ServicePhoto 
        WHERE
            serviceId = ?
    `

	photos := make([]types.ServicePhoto, 0)
	err := db.Connection.Select(&photos, query, serviceId)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
