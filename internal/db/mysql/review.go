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

	tx, err := m.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	insertReview := `
        INSERT INTO
            Review (serviceId, rating)
        VALUES
            (:serviceId, :rating)
    `

	res, err := tx.NamedExecContext(ctx, insertReview, review)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	updateReviewedFlag := `
        UPDATE 
            Transaction 
        SET 
            isReviewed = 1
        WHERE
            id = ?
    `
	_, err = tx.ExecContext(ctx, updateReviewedFlag, review.TransactionId)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(insertId), nil
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
