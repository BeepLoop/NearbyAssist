package review

import (
	"nearbyassist/internal/models"

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
	return "", nil
}

func (s *MysqlReviewStore) FindById(id string) (*models.ReviewModel, error) {
	return nil, nil
}

func (s *MysqlReviewStore) FindAll() ([]*models.ReviewModel, error) {
	return nil, nil
}

// Returns nil if reviewable, else error
func (s *MysqlReviewStore) IsServiceReviewable(serviceId string) error {
	return nil
}

func (s *MysqlReviewStore) GetTransactionById(id string) (*models.TransactionModel, error) {
	return nil, nil
}

func (s *MysqlReviewStore) GetReviewsByService(serviceId string) ([]*models.ReviewModel, error) {
	return nil, nil
}
