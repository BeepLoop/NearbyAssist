package review

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlReviewStore struct {
	db *sqlx.DB
}

func NewMysqlReviewStore(db *sqlx.DB) *MysqlReviewStore {
	return &MysqlReviewStore{
		db: db,
	}
}

func (s *MysqlReviewStore) Create(data *models.ReviewModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	insertReview := "INSERT INTO Review (id, serviceId, rating) VALUES (:id, :serviceId, :rating)"
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

func (s *MysqlReviewStore) FindById(id string) (*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	review := new(models.ReviewModel)
	query := "SELECT id, serviceId, rating FROM Review WHERE id = ?"
	if err := s.db.GetContext(ctx, review, query, id); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *MysqlReviewStore) IsServiceReviewable(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	isReviewed := true
	query := "SELECT isReviewed FROM Transaction WHERE serviceId = ?"
	if err := s.db.GetContext(ctx, isReviewed, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if isReviewed {
		return errors.New("service is already reviewed")
	}

	return nil
}

func (s *MysqlReviewStore) GetTransactionById(id string) (*models.TransactionModel, error) {
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

func (s *MysqlReviewStore) GetReviewsByService(serviceId string) ([]*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]*models.ReviewModel, 0)
	if err := s.db.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
	}

	return reviews, nil
}
