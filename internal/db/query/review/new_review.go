package review_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func CreateReview(review types.Review) error {
	query := `
        INSERT INTO
            Review (vendorId, rating) 
        VALUES 
            (:vendorId, :rating)
    `

	_, err := db.Connection.NamedExec(query, review)
	if err != nil {
		return err
	}

	return nil
}
