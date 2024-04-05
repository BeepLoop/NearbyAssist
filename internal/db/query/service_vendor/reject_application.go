package vendor_query

import "nearbyassist/internal/db"

func RejectApplication(applicationId int) error {
	query := `
        UPDATE
            Application 
        SET 
            status = 'rejected'
        WHERE 
            id = ?
    `

	_, err := db.Connection.Exec(query, applicationId)
	if err != nil {
		return err
	}

	return nil
}
