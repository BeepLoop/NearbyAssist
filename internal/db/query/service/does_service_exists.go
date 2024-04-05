package service_query

import (
	"nearbyassist/internal/db"
)

func DoesServiceExists(serviceId, vendorId int) bool {
	query := `
        SELECT
            id 
        FROM 
            Service 
        WHERE
            id = ? AND vendorId = ?
    `

	result := new(struct {
		Id int `db:"id"`
	})
	err := db.Connection.Get(result, query, serviceId, vendorId)
	if err != nil {
		return false
	}

	return true
}
