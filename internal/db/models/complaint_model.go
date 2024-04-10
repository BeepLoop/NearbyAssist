package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

type ComplaintModel struct {
	Model
	UpdateableModel
	Code    int    `json:"code" db:"code"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

func NewComplaintModel(db *db.DB) *ComplaintModel {
	return &ComplaintModel{
		Model: Model{Db: db},
	}
}

func (c *ComplaintModel) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            Complaint (vendorId, code, title, content)
        VALUES
            (:vendorId, :code, :title, :content)
    `

	res, err := c.Db.Conn.NamedExecContext(ctx, query, c)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (c *ComplaintModel) Update(id int) error {
	return nil
}

func (c *ComplaintModel) Delete(id int) error {
	return nil
}

func (c *ComplaintModel) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            COUNT(*)
        FROM
            Complaint
    `

	count := 0
	err := c.Db.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return 0, nil
}
