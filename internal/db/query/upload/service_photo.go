package upload_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func UploadServicePhoto(data types.UploadData) error {
	query := `
        INSERT INTO
            ServicePhoto (vendorId, serviceId, url)
        VALUES 
            (:vendorId, :serviceId, :imageUrl)
    `

	_, err := db.Connection.NamedExec(query, data)
	if err != nil {
		return err
	}

	return nil
}
