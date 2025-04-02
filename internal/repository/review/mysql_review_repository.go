package review_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
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

	data.Id = utils.GenerateId()

	insertReview := "INSERT INTO Review (id, transactionId, rating, text) VALUES (:id, :transactionId, :rating, :text)"
	if _, err := tx.NamedExecContext(ctx, insertReview, data); err != nil {
		return "", err
	}

	updateReviewedFlag := "UPDATE Transaction SET isReviewed = 1 WHERE id = ?"
	if _, err = tx.ExecContext(ctx, updateReviewedFlag, data.TransactionId); err != nil {
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
	query := "SELECT id, transactionId, rating FROM Review WHERE id = ?"
	if err := s.db.GetContext(ctx, review, query, id); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *MysqlReviewRepository) FindByTransactionId(id string) (*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transaction := new(models.TransactionModel)
	query := "SELECT * FROM Transaction WHERE id = ?"
	if err := s.db.GetContext(ctx, transaction, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transaction, nil
}
