package complaint_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func CreateComplaint(complaint types.Complaint) error {
	query := `
        INSERT INTO 
            Complaint (vendorId, code, title, content)
        VALUES
            (:vendorId, :code, :title, :content)
    `

	_, err := db.Connection.NamedExec(query, complaint)
	if err != nil {
		return err
	}

	return nil
}
