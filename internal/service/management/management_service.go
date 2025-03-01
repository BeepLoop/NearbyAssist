package management_service

import (
	"nearbyassist/internal/models"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/service/auth"
)

type Service struct {
	store   user_repo.UserRepository
	encrypt auth.Encryption
	hash    auth.Hash
}

func NewService(store user_repo.UserRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{store: store, encrypt: encrypt, hash: hash}
}

func (s *Service) GetUsers(limit, offset int) ([]*models.UserModel, error) {
	accounts, err := s.store.GetAllUserAccounts(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if decrypted, err := s.encrypt.DecryptString(account.Name); err != nil {
			return nil, err
		} else {
			account.Name = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(account.Email); err != nil {
			return nil, err
		} else {
			account.Email = decrypted
		}
	}

	return accounts, nil
}

func (s *Service) FindUserByEmail(email string) (*models.UserModel, error) {
	emailHash, err := s.hash.Generate([]byte(email))
	if err != nil {
		return nil, err
	}

	user, err := s.store.FindByEmailHash(emailHash)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(user.Name); err != nil {
		return nil, err
	} else {
		user.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(user.Email); err != nil {
		return nil, err
	} else {
		user.Email = decrypted
	}

	return user, nil
}

func (s *Service) GetSingleUser(userId string) (*models.UserAccountPageData, error) {
	accountData, err := s.store.GetUserAccountPageData(userId)
	if err != nil {
		return nil, err
	}

	if stat, err := s.store.GetSentTransactionCount(userId); err != nil {
		accountData.Stat.Sent = models.SentStat{}
	} else {
		accountData.Stat.Sent = *stat
	}

	if stat, err := s.store.GetReceivedTransactionCount(userId); err != nil {
		accountData.Stat.Received = models.ReceivedStat{}
	} else {
		accountData.Stat.Received = *stat
	}

	if decrypted, err := s.encrypt.DecryptString(accountData.Name); err != nil {
		return nil, err
	} else {
		accountData.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(accountData.Email); err != nil {
		return nil, err
	} else {
		accountData.Email = decrypted
	}

	if accountData.Address.Valid {
		if decrypted, err := s.encrypt.DecryptString(accountData.Address.String); err != nil {
			return nil, err
		} else {
			accountData.Address.String = decrypted
		}
	} else {
		accountData.Address.String = ""
	}

	for _, service := range accountData.Services {
		if decrypted, err := s.encrypt.DecryptString(service.Title); err != nil {
			return nil, err
		} else {
			service.Title = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(service.Description); err != nil {
			return nil, err
		} else {
			service.Description = decrypted
		}
	}

	return accountData, nil
}

func (s *Service) BanUser(userId string) error {
	if err := s.store.BanUser(userId); err != nil {
		return err
	}

	return nil
}

func (s *Service) UnbanUser(userId string) error {
	if err := s.store.UnbanUser(userId); err != nil {
		return err
	}

	return nil
}
