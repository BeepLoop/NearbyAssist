package upload_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func UploadServicePhoto(data types.ServicePhoto) (int, error) {
	query := `
        INSERT INTO
            ServicePhoto (vendorId, serviceId, url)
        VALUES 
            (:vendorId, :serviceId, :url)
    `

	res, err := db.Connection.NamedExec(query, data)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
