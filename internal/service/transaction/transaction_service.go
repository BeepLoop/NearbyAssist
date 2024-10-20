package transaction_service

import (
	"nearbyassist/internal/models"
	transaction_repo "nearbyassist/internal/repository/transaction"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   transaction_repo.TransactionRepository
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store transaction_repo.TransactionRepository, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{store: store, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateTransaction(req *request.NewTransactionPayload) (string, error) {
	// Validate that the date is valid
	if err := utils.ValidateDateRange(req.Start, req.End); err != nil {
		if err.Error() == utils.DATE_PARSE_ERR {
			return "", err
		}

		return "", err
	}

	transaction := new(models.TransactionModel)
	transaction.ClientId = req.ClientId
	transaction.VendorId = req.VendorId
	transaction.ServiceId = req.ServiceId
	transaction.Start = req.Start
	transaction.End = req.End

	transactionId, err := s.store.Create(transaction)
	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (s *Service) GetUserTransactionList(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetMyTransactions(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		if plain, err := s.encrypt.DecryptString(transaction.Vendor); err != nil {
			return nil, err
		} else {
			transaction.Vendor = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Client); err != nil {
			return nil, err
		} else {
			transaction.Client = plain
		}
	}

	return transactions, nil
}

func (s *Service) GetOngoingTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetOngoing(userId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) GetTransactionHistory(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetHistory(userId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) CompleteTransaction(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if transaction, err := s.store.FindById(transactionId); err != nil {
		return err
	} else {
		if transaction.ClientId != userId {
			return err
		}
	}

	if err := s.store.MarkComplete(transactionId); err != nil {
		return err
	}

	return nil
}
