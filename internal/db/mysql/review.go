package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"time"
)

func (m *Mysql) CreateReview(review *request.NewReview) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Review (serviceId, rating)
        VALUES
            (:serviceId, :rating)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, review)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) FindReviewById(id int) (*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE id = ?"

	review := models.NewReviewModel()
	err := m.Conn.GetContext(ctx, review, query, id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (m *Mysql) FindAllReviewByService(id int) ([]models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]models.ReviewModel, 0)
	err := m.Conn.SelectContext(ctx, &reviews, query, id)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
