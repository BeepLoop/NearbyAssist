package transaction_service

import (
	"errors"
	"nearbyassist/internal/models"
	transaction_repo "nearbyassist/internal/repository/transaction"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
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
	transaction := new(models.TransactionModel)
	transaction.ClientId = req.ClientId
	transaction.VendorId = req.VendorId
	transaction.ServiceId = req.ServiceId
	transaction.Cost = req.Cost

	extras := make([]models.ExtraModel, 0)
	for _, extra := range req.Extras {
		extras = append(extras, models.ExtraModel{
			Model: models.Model{
				Id: extra.Id,
			},
			Title:       extra.Title,
			Description: extra.Description,
			Price:       extra.Price,
		})
	}

	transaction.Extras = extras

	transactionId, err := s.store.Create(transaction)
	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (s *Service) GetTransaction(transactionId string) (*models.TransactionModel, error) {
	transaction, err := s.store.FindById(transactionId)
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

	extras := make([]models.ExtraModel, 0)
	for _, extra := range transaction.Extras {
		title, err := s.encrypt.DecryptString(extra.Title)
		if err != nil {
			return nil, err
		}

		description, err := s.encrypt.DecryptString(extra.Description)
		if err != nil {
			return nil, err
		}

		extras = append(extras, models.ExtraModel{
			Model:       extra.Model,
			Title:       title,
			Description: description,
			Price:       extra.Price,
		})
	}
	transaction.Extras = extras

	return transaction, nil
}

func (s *Service) CancelTransaction(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.store.FindById(transactionId)
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

	if err := s.store.Cancel(transactionId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransactionSummary(transactionId string) (*response.TransactionSummary, error) {
	transactionData, err := s.store.GetSummary(transactionId)
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

func (s *Service) GetTransactionUserSent(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetTransactionSent(userId)
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

		extras := make([]models.ExtraModel, 0)
		for _, extra := range transaction.Extras {
			title, err := s.encrypt.DecryptString(extra.Title)
			if err != nil {
				return nil, err
			}

			description, err := s.encrypt.DecryptString(extra.Description)
			if err != nil {
				return nil, err
			}

			extras = append(extras, models.ExtraModel{
				Model:       extra.Model,
				Title:       title,
				Description: description,
				Price:       extra.Price,
			})
		}
		transaction.Extras = extras
	}

	return transactions, nil
}

func (s *Service) GetTransactionUserReceived(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetTransactionReceived(userId)
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

		extras := make([]models.ExtraModel, 0)
		for _, extra := range transaction.Extras {
			title, err := s.encrypt.DecryptString(extra.Title)
			if err != nil {
				return nil, err
			}

			description, err := s.encrypt.DecryptString(extra.Description)
			if err != nil {
				return nil, err
			}

			extras = append(extras, models.ExtraModel{
				Model:       extra.Model,
				Title:       title,
				Description: description,
				Price:       extra.Price,
			})
		}
		transaction.Extras = extras
	}

	return transactions, nil
}

func (s *Service) GetRecentTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetRecent(userId)
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

		extras := make([]models.ExtraModel, 0)
		for _, extra := range transaction.Extras {
			title, err := s.encrypt.DecryptString(extra.Title)
			if err != nil {
				return nil, err
			}

			description, err := s.encrypt.DecryptString(extra.Description)
			if err != nil {
				return nil, err
			}

			extras = append(extras, models.ExtraModel{
				Model:       extra.Model,
				Title:       title,
				Description: description,
				Price:       extra.Price,
			})
		}
		transaction.Extras = extras
	}

	return transactions, nil
}

func (s *Service) GetConfirmedTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetConfirmed(userId)
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

		extras := make([]models.ExtraModel, 0)
		for _, extra := range transaction.Extras {
			title, err := s.encrypt.DecryptString(extra.Title)
			if err != nil {
				return nil, err
			}

			description, err := s.encrypt.DecryptString(extra.Description)
			if err != nil {
				return nil, err
			}

			extras = append(extras, models.ExtraModel{
				Model:       extra.Model,
				Title:       title,
				Description: description,
				Price:       extra.Price,
			})
		}
		transaction.Extras = extras
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

		extras := make([]models.ExtraModel, 0)
		for _, extra := range transaction.Extras {
			title, err := s.encrypt.DecryptString(extra.Title)
			if err != nil {
				return nil, err
			}

			description, err := s.encrypt.DecryptString(extra.Description)
			if err != nil {
				return nil, err
			}

			extras = append(extras, models.ExtraModel{
				Model:       extra.Model,
				Title:       title,
				Description: description,
				Price:       extra.Price,
			})
		}
		transaction.Extras = extras
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
