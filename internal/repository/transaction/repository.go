package transaction_repo

import "nearbyassist/internal/models"

type TransactionRepository interface {
	Create(data *models.TransactionModel) (string, error)
	FindById(id string) (*models.TransactionModel, error)
	GetAll() ([]*models.TransactionModel, error)

	GetMyTransactions(id string) ([]*models.TransactionModel, error)

	GetOngoing(id string) ([]*models.TransactionModel, error)

	GetHistory(id string) ([]*models.TransactionModel, error)

	MarkComplete(transactionId string) error
}
