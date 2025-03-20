package transaction_service

import (
	"errors"
	"nearbyassist/internal/models"
	transaction_repo "nearbyassist/internal/repository/transaction"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	transactionStore transaction_repo.TransactionRepository
	encrypt          core.Encryption
	jwt              core.Authenticator
}

func NewService(transactionStore transaction_repo.TransactionRepository, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{
		transactionStore: transactionStore,
		encrypt:          encrypt,
		jwt:              jwt,
	}
}

func (s *Service) CreateTransaction(req *request.NewTransactionPayload) (string, error) {
	transaction := new(models.TransactionModel)
	transaction.ClientId = req.ClientId
	transaction.VendorId = req.VendorId
	transaction.ServiceId = req.ServiceId
	transaction.Cost = req.Cost

	extras := make([]*models.ExtraModel, 0)
	for _, extra := range req.Extras {
		extras = append(extras, &models.ExtraModel{
			Model: models.Model{
				Id: extra.Id,
			},
		})
	}

	transaction.Extras = extras

	transactionId, err := s.transactionStore.Create(transaction)
	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (s *Service) GetTransaction(transactionId string) (*models.TransactionModel, error) {
	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return nil, err
	}

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

	if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
		return nil, err
	} else {
		transaction.Service.Title = plain
	}

	if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
		return nil, err
	} else {
		transaction.Service.Description = plain
	}

	for _, extra := range transaction.Extras {
		if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
			return nil, err
		} else {
			extra.Title = title
		}

		if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
			return nil, err
		} else {
			extra.Description = description
		}
	}

	return transaction, nil
}

func (s *Service) CancelTransaction(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return err
	}

	if transaction.Status != models.TRANSACTION_STATUS_PENDING {
		return errors.New("Could not cancel non-pending transaction")
	}

	if transaction.Status == models.TRANSACTION_STATUS_DONE || transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		return errors.New("Transaction already completed or cancelled")
	}

	if transaction.ClientId != userId {
		return errors.New("Unauthorized cancel request")
	}

	if err := s.transactionStore.Cancel(transactionId); err != nil {
		return err
	}

	return nil
}

func (s *Service) AcceptTransactionRequest(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return err
	}

	if transaction.Status != models.TRANSACTION_STATUS_PENDING {
		return errors.New("Could not accept non-pending transaction")
	}

	if transaction.Status == models.TRANSACTION_STATUS_DONE || transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		return errors.New("Transaction already completed or cancelled")
	}

	if transaction.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	if err := s.transactionStore.Accept(transactionId); err != nil {
		return err
	}

	return nil
}

func (s *Service) RejectTransactionRequest(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return err
	}

	if transaction.Status != models.TRANSACTION_STATUS_PENDING {
		return errors.New("Could not reject non-pending transaction")
	}

	if transaction.Status == models.TRANSACTION_STATUS_DONE || transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		return errors.New("Transaction already completed or cancelled")
	}

	if transaction.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	if err := s.transactionStore.Reject(transactionId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransactionSummary(transactionId string) (*response.TransactionSummary, error) {
	transactionData, err := s.transactionStore.GetSummary(transactionId)
	if err != nil {
		return nil, err
	}

	if plain, err := s.encrypt.DecryptString(transactionData.Vendor); err != nil {
		return nil, err
	} else {
		transactionData.Vendor = plain
	}

	if plain, err := s.encrypt.DecryptString(transactionData.Client); err != nil {
		return nil, err
	} else {
		transactionData.Client = plain
	}

	if plain, err := s.encrypt.DecryptString(transactionData.ServiceTitle); err != nil {
		return nil, err
	} else {
		transactionData.ServiceTitle = plain
	}

	if plain, err := s.encrypt.DecryptString(transactionData.VendorEmail); err != nil {
		return nil, err
	} else {
		transactionData.VendorEmail = plain
	}

	if plain, err := s.encrypt.DecryptString(transactionData.ClientEmail); err != nil {
		return nil, err
	} else {
		transactionData.ClientEmail = plain
	}

	return transactionData, nil
}

func (s *Service) GetUserTransactionList(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetMyTransactions(userId)
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

func (s *Service) GetTransactionUserSent(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetTransactionSent(userId)
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

		if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
			return nil, err
		} else {
			transaction.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
			return nil, err
		} else {
			transaction.Service.Description = plain
		}

		for _, extra := range transaction.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}
		}
	}

	return transactions, nil
}

func (s *Service) GetTransactionUserReceived(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetTransactionReceived(userId)
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

		if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
			return nil, err
		} else {
			transaction.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
			return nil, err
		} else {
			transaction.Service.Description = plain
		}

		for _, extra := range transaction.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}
		}
	}

	return transactions, nil
}

func (s *Service) GetRecentTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetRecent(userId)
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

		if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
			return nil, err
		} else {
			transaction.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
			return nil, err
		} else {
			transaction.Service.Description = plain
		}

		for _, extra := range transaction.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}
		}
	}

	return transactions, nil
}

func (s *Service) GetConfirmedTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetConfirmed(userId)
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

		if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
			return nil, err
		} else {
			transaction.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
			return nil, err
		} else {
			transaction.Service.Description = plain
		}

		for _, extra := range transaction.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}
		}
	}

	return transactions, nil
}

func (s *Service) GetReviewableTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	reviewables, err := s.transactionStore.GetReviewableTransactions(userId)
	if err != nil {
		return nil, err
	}

	for _, reviewable := range reviewables {
		if plain, err := s.encrypt.DecryptString(reviewable.Vendor); err != nil {
			return nil, err
		} else {
			reviewable.Vendor = plain
		}

		if plain, err := s.encrypt.DecryptString(reviewable.Client); err != nil {
			return nil, err
		} else {
			reviewable.Client = plain
		}

		if plain, err := s.encrypt.DecryptString(reviewable.Service.Title); err != nil {
			return nil, err
		} else {
			reviewable.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(reviewable.Service.Description); err != nil {
			return nil, err
		} else {
			reviewable.Service.Description = plain
		}

		for _, extra := range reviewable.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}
		}
	}

	return reviewables, nil
}

func (s *Service) GetTransactionHistory(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetHistory(userId)
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

		if plain, err := s.encrypt.DecryptString(transaction.Service.Title); err != nil {
			return nil, err
		} else {
			transaction.Service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(transaction.Service.Description); err != nil {
			return nil, err
		} else {
			transaction.Service.Description = plain
		}

		for _, extra := range transaction.Extras {
			if title, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = title
			}

			if description, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = description
			}

		}
	}

	return transactions, nil
}

func (s *Service) CompleteTransaction(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if transaction, err := s.transactionStore.FindById(transactionId); err != nil {
		return err
	} else {
		if transaction.VendorId != userId {
			return errors.New("unauthorized")
		}
	}

	if err := s.transactionStore.MarkComplete(transactionId); err != nil {
		return err
	}

	return nil
}
