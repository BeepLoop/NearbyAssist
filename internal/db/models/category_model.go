package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

type CategoryModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`
}

func NewCategoryModel(db *db.DB) *CategoryModel {
	return &CategoryModel{
		Model: Model{Db: db},
	}
}

func (c *CategoryModel) Create() (int, error) {
	return 0, nil
}

func (c *CategoryModel) Update(id int) error {
	return nil
}

func (c *CategoryModel) Delete(id int) error {
	return nil
}

func (c *CategoryModel) FindAll() ([]*CategoryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, title
        FROM 
            Category
    `

	categories := make([]*CategoryModel, 0)
	err := c.Db.Conn.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return categories, nil
}
