package vendor_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetApplicants(filter string) ([]types.Application, error) {
	query := QueryFactory(filter)

	applicants := make([]types.Application, 0)
	err := db.Connection.Select(&applicants, query)
	if err != nil {
		return nil, err
	}

	return applicants, nil
}

func QueryFactory(filter string) string {
	defaultQeury := `
        SELECT 
            id, applicantId, job, status
        FROM
            Application
    `

	switch filter {
	case "":
		return defaultQeury
	case "pending":
		return `
            SELECT 
                id, applicantId, job, status
            FROM
                Application
            WHERE
                status = 'pending'
        `
	case "approved":
		return `
            SELECT 
                id, applicantId, job, status
            FROM
                Application
            WHERE
                status = 'approved'
        `
	case "rejected":
		return `
            SELECT 
                id, applicantId, job, status
            FROM
                Application
            WHERE
                status = 'rejected'
        `
	default:
		return defaultQeury
	}
}
