package review

import "nearbyassist/internal/models"

type ReviewStore interface {
	Create(data *models.ReviewModel) (string, error)
	FindById(id string) (*models.ReviewModel, error)
	FindAll() ([]*models.ReviewModel, error)

	// Returns nil if reviewable, else error
	IsServiceReviewable(serviceId string) error

	GetTransactionById(id string) (*models.TransactionModel, error)

	GetReviewsByService(serviceId string) ([]*models.ReviewModel, error)
}
