package query

import (
	"fmt"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func RegisterLocation(location types.LocationRegister) error {
	query := `
	        INSERT INTO
	            Location (ownerId, address, location)
	        VALUES 
	            (1, ?, ST_GeomFromText(?, 4326))
	    `

	point := fmt.Sprintf("POINT(%f %f)", location.Latitude, location.Longitude)

	_, err := db.Connection.Exec(query, location.Address, point)
	if err != nil {
		return err
	}

	return nil
}