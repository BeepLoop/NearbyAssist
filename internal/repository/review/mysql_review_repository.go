package review_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlReviewRepository struct {
	db *sqlx.DB
}

func NewMysqlReviewRepository(db *sqlx.DB) *MysqlReviewRepository {
	return &MysqlReviewRepository{
		db: db,
	}
}

func (s *MysqlReviewRepository) Create(data *models.ReviewModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertReview := `
        INSERT INTO
            Review (revieweeId, bookingId, rating, text)
        VALUES
            (:revieweeId, :bookingId, :rating, :text)
    `
	if _, err := tx.NamedExecContext(ctx, insertReview, data); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlReviewRepository) FindById(id string) (*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	review := new(models.ReviewModel)
	query := "SELECT id, bookingId, rating FROM Review WHERE id = ?"
	if err := s.db.GetContext(ctx, review, query, id); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *MysqlReviewRepository) GetUserReviewOnBooking(userId, bookingId string) (*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            *
        FROM
            Review
        WHERE
            revieweeId = ? AND bookingId = ?
    `

	review := new(models.ReviewModel)
	if err := s.db.GetContext(ctx, review, query, userId, bookingId); err != nil {
		return nil, err
	}

	getRevieweeQuery := `
        SELECT
            id, name, email, imageUrl
        FROM
            User
        WHERE
            id = ?
    `
	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, getRevieweeQuery, userId); err != nil {
		return nil, err
	}
	review.Reviewee = user

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return review, nil
}
