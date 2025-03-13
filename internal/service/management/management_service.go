package management_service

import (
	"fmt"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/service/auth"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/utils"
	"time"
)

type Service struct {
	store      user_repo.UserRepository
	notifStore notification_repo.NotificationRepository
	encrypt    auth.Encryption
	hash       auth.Hash
}

func NewService(store user_repo.UserRepository, notifStore notification_repo.NotificationRepository, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{
		store:      store,
		notifStore: notifStore,
		encrypt:    encrypt,
		hash:       hash,
	}
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

	if restricted, expired, err := s.store.IsRestricted(userId); err != nil {
		return nil, err
	} else {
		accountData.Restricted = restricted && !expired
	}

	if accountData.Restricted {
		if err := s.store.LiftRestrictionIfExpired(userId); err != nil {
			return nil, err
		}
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
		// Decrypt service title and description
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

		// Decrypt service extra title and description
		for _, extra := range service.Extras {
			if decrypted, err := s.encrypt.DecryptString(extra.Title); err != nil {
				return nil, err
			} else {
				extra.Title = decrypted
			}

			if decrypted, err := s.encrypt.DecryptString(extra.Description); err != nil {
				return nil, err
			} else {
				extra.Description = decrypted
			}
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

func (s *Service) RestrictUser(userId, reason, duration string) error {
	d, err := utils.ParseStringDuration(duration)
	if err != nil {
		return err
	}

	endDate := time.Now().Add(d)
	data := &models.RestrictionModel{
		UserId:  userId,
		Reason:  reason,
		EndTime: utils.FormatDateTime(endDate),
	}

	if err := s.store.RestrictUser(data); err != nil {
		return err
	}

	stringifiedDuration := utils.FormatDurationToString(d)

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     "Account Restricted " + stringifiedDuration,
		Content:   reason,
	}

	if encrypted, err := s.encrypt.EncryptString(notification.Title); err != nil {
		return err
	} else {
		notification.Title = encrypted
	}

	if encrypted, err := s.encrypt.EncryptString(notification.Content); err != nil {
		return err
	} else {
		notification.Content = encrypted
	}

	if err := s.notifStore.Create(notification); err != nil {
		return err
	}

	notificationHeading := "Account Restricted!"
	notificationContent := "You commited a violation resulting to account restriction."

	oneSignal := notification_service.OneSignalInstance
	if oneSignal != nil {
		if err := oneSignal.NewUrgentNotification(userId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	return nil
}

func (s *Service) UnrestrictUser(userId string) error {
	if err := s.store.ForceLiftRestriction(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     "Restriction Lifted",
		Content:   "The restriction to your account has been lifted by the administrator. Avoid committing violations to prevent future restrictions.",
	}

	if encrypted, err := s.encrypt.EncryptString(notification.Title); err != nil {
		return err
	} else {
		notification.Title = encrypted
	}

	if encrypted, err := s.encrypt.EncryptString(notification.Content); err != nil {
		return err
	} else {
		notification.Content = encrypted
	}

	if err := s.notifStore.Create(notification); err != nil {
		return err
	}

	notificationHeading := "Account Restriction Lifted!"
	notificationContent := "Your account restriction has been lifted."

	oneSignal := notification_service.OneSignalInstance
	if oneSignal != nil {
		if err := oneSignal.NewUrgentNotification(userId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	return nil
}
