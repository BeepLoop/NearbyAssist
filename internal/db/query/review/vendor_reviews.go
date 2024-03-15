package review_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func VendorReview(vendorId string) ([]types.Review, error) {
	query := `
        SELECT 
            id, vendorId, rating
        FROM 
            Review 
        WHERE
            vendorId = ?
    `

	reviews := make([]types.Review, 0)
	err := db.Connection.Select(&reviews, query, vendorId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
