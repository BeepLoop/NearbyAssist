package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindAllCategory() ([]models.CategoryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, title FROM Category"

	categories := make([]models.CategoryModel, 0)
	err := m.Conn.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return categories, nil
}
