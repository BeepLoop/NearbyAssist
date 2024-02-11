package db

import (
	"fmt"
	"nearbyassist/internal/types"
)

func RegisterLocation(location types.LocationRegister) error {
	query := `
	        INSERT INTO
	            Location (address, location)
	        VALUES 
	            (?, ST_GeomFromText(?, 4326))
	    `

	point := fmt.Sprintf("POINT(%f %f)", location.Latitude, location.Longitude)

	_, err := DB_CONN.Exec(query, location.Address, point)
	if err != nil {
		return err
	}

	return nil
}
