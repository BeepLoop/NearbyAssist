package review_repo

import "nearbyassist/internal/models"

type ReviewRepository interface {
	Create(data *models.ReviewModel) (string, error)
	FindById(id string) (*models.ReviewModel, error)

	IsReviewed(serviceId string) (bool, error)

	GetTransactionById(id string) (*models.TransactionModel, error)

	GetReviewsByService(serviceId string) ([]*models.ReviewModel, error)
}
