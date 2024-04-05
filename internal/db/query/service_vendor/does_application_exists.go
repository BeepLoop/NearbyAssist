package vendor_query

import "nearbyassist/internal/db"

func DoesApplicationExists(applicationId, applicantId int) bool {
	query := `
        SELECT
            id 
        FROM 
            Application 
        WHERE
            id = ? AND applicantId = ?
    `

	result := new(struct {
		Id int `db:"id"`
	})
	err := db.Connection.Get(result, query, applicationId, applicantId)
	if err != nil {
		return false
	}

	return true
}
