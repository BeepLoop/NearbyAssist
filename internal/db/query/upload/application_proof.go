package upload_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func UploadApplicationProof(data types.ApplicationProof) (int, error) {
	query := `
        INSERT INTO
            ApplicationProof (applicationId, applicantId, url)
        VALUES
            (:applicationId, :applicantId, :url)
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
