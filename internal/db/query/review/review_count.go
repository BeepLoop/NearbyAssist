package review_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func ReviewCount(vendorId string) (types.Count, error) {
	query := `
        SELECT 
            rating, COUNT(*) AS count 
        FROM 
            Review 
        WHERE 
            vendorId = ?
        GROUP BY 
            rating;
    `

	reviewCount := make([]types.ReviewCount, 0)
	err := db.Connection.Select(&reviewCount, query, vendorId)
	if err != nil {
		return nil, err
	}

	countMap := types.InitCountMap()
	for _, count := range reviewCount {
		countMap[count.Rating] = count.Count
	}

	return countMap, nil
}
