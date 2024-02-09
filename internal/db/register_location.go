package db

import "nearbyassist/internal/types"

func RegisterLocation(location types.Location) error {
	query := `
        INSERT INTO
            Location (address, latitude, longitude)
        VALUES 
            (?, ?, ?)
    `

	_, err := DB_CONN.Exec(query, location.Address, location.Latitude, location.Longitude)
	if err != nil {
		return err
	}

	return nil
}
