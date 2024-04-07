package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	Rating    string `json:"rating" db:"rating"`
}

func NewReviewModel() *ReviewModel {
	return &ReviewModel{}
}

func (r *ReviewModel) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Review (serviceId, rating) 
        VALUES 
            (:serviceId, :rating)
    `

	res, err := db.Connection.NamedExecContext(ctx, query, r)
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

func (r *ReviewModel) Update(id int) error {
	return nil
}

func (r *ReviewModel) Delete(id int) error {
	return nil
}

func (r *ReviewModel) FindById(id int) (*ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, serviceId, rating
        FROM
            Review
        WHERE
            id = ?
    `

	review := new(ReviewModel)
	err := db.Connection.GetContext(ctx, review, query, id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewModel) FindByService(serviceId int) ([]*ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, serviceId, rating
        FROM
            Review
        WHERE
            serviceId = ?
    `

	reviews := make([]*ReviewModel, 0)
	err := db.Connection.SelectContext(ctx, &reviews, query, serviceId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
