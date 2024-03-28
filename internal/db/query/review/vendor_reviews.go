package review_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func VendorReview(vendorId int) ([]types.Review, error) {
	query := `
        SELECT 
            id, serviceId, rating
        FROM 
            Review 
        WHERE
            serviceId = ?
    `

	reviews := make([]types.Review, 0)
	err := db.Connection.Select(&reviews, query, vendorId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
