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

	GetOngoing(id string) ([]*models.TransactionModel, error)

	GetHistory(id string) ([]*models.TransactionModel, error)

	MarkComplete(transactionId string) error
}
