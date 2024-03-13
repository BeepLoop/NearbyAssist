package photo_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func UploadPhoto(data types.UploadData) error {
	query := `
        INSERT INTO
            Photo (vendor, service, url)
        VALUES 
            (:vendorId, :serviceId, :imageUrl)
    `

	_, err := db.Connection.NamedExec(query, data)
	if err != nil {
		return err
	}

	return nil
}
