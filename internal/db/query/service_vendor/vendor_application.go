package vendor_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func VendorApplication(data types.VendorApplication) (int, error) {
	query := `
        INSERT INTO 
            Application (applicantId, job, latitude, longitude)
        VALUES
            (:applicantId, :job, :latitude, :longitude)
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
