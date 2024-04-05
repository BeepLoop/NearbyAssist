package vendor_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetAllApplicants() ([]types.Application, error) {
	query := `
        SELECT 
            id, applicantId, job, status
        FROM
            Application
    `

	applicants := make([]types.Application, 0)
	err := db.Connection.Select(&applicants, query)
	if err != nil {
		return nil, err
	}

	return applicants, nil
}
