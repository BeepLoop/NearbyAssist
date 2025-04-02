package transaction_repo

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
)

type TransactionRepository interface {
	Create(data *models.TransactionModel) (string, error)
	FindById(id string) (*models.TransactionModel, error)
	GetSummary(id string) (*response.TransactionSummary, error)
	GetAll() ([]*models.TransactionModel, error)

	GetMyTransactions(id string) ([]*models.TransactionModel, error)
	GetTransactionSent(id string) ([]*models.TransactionModel, error)
	GetTransactionReceived(id string) ([]*models.TransactionModel, error)

	GetRecent(id string) ([]*models.TransactionModel, error)
	GetConfirmed(id string) ([]*models.TransactionModel, error)

	GetHistory(id string) ([]*models.TransactionModel, error)
	GetReviewableTransactions(userId string) ([]*models.TransactionModel, error)

	Cancel(transactionId string) error
	Accept(transactionId string) error
	Reject(transactionId string) error
	MarkComplete(transactionId string) error

	IsReviewed(transactionId string) (bool, error)
}
