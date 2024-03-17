package service_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func RegisterService(service types.ServiceRegister) error {
	query := `
	        INSERT INTO
	            Service
	                (vendorId, title, description, rate, location, category)
	        VALUES 
                (
                    :vendorId,
                    :title,
                    :description,
                    :rate,
                    ST_GeomFromText(:point, 4326),
                    :categoryId
                )
	    `

	_, err := db.Connection.NamedExec(query, service)
	if err != nil {
		return err
	}

	return nil
}
