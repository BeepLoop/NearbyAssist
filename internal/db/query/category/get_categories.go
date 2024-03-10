package category_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func GetCategories() ([]types.Category, error) {
	query := `
        SELECT 
            id, title
        FROM 
            Category
    `

	categories := make([]types.Category, 0)
	err := db.Connection.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
