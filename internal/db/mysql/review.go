package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"time"
)

func (m *Mysql) CreateReview(review *request.NewReview) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := m.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertReview := `
        INSERT INTO
            Review (id, serviceId, rating)
        VALUES
            (:id, :serviceId, :rating)
    `

	if _, err := tx.NamedExecContext(ctx, insertReview, review); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
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
			return "", err
		}

		return "", err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return review.Id, nil
}

func (m *Mysql) FindReviewById(reviewId string) (*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE id = ?"

	review := models.NewReviewModel()
	if err := m.Conn.GetContext(ctx, review, query, reviewId); err != nil {
		return nil, err
	}

	return review, nil
}

func (m *Mysql) FindAllReviewByService(serviceId string) ([]models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]models.ReviewModel, 0)
	if err := m.Conn.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
	}

	return reviews, nil
}
