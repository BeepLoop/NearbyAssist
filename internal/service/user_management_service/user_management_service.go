package user_management_service

import (
	"fmt"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"time"
)

type Service struct {
	userStore  user_repo.UserRepository
	notifStore notification_repo.NotificationRepository
	ws         websocket.Socket
	encrypt    core.Encryption
	hash       core.Hash
}

func NewService(userStore user_repo.UserRepository, notifStore notification_repo.NotificationRepository, ws websocket.Socket, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		userStore:  userStore,
		notifStore: notifStore,
		ws:         ws,
		encrypt:    encrypt,
		hash:       hash,
	}
}

func (s *Service) GetSingleUser(userId string) (*models.UserAccountPageData, error) {
	accountData, err := s.userStore.GetUserAccountPageData(userId)
	if err != nil {
		return nil, err
	}

	if restricted, expired, err := s.userStore.IsRestricted(userId); err != nil {
		return nil, err
	} else {
		accountData.Restricted = restricted && !expired
	}

	if accountData.Restricted {
		if err := s.userStore.LiftRestrictionIfExpired(userId); err != nil {
			return nil, err
		}
	}

	if stat, err := s.userStore.GetSentTransactionCount(userId); err != nil {
		accountData.Stat.Sent = models.SentStat{}
	} else {
		accountData.Stat.Sent = *stat
	}

	if stat, err := s.userStore.GetReceivedTransactionCount(userId); err != nil {
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
	if err := s.userStore.BanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

	return nil
}

func (s *Service) UnbanUser(userId string) error {
	if err := s.userStore.UnbanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

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

	if err := s.userStore.RestrictUser(data); err != nil {
		return err
	}

	stringifiedDuration := utils.FormatDurationToString(d)

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     "Account Restricted " + stringifiedDuration,
		Content:   reason,
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	}

	notificationHeading := "Account Restricted!"
	notificationContent := "You commited a violation resulting to account restriction."

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(userId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: userId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	syncEvent := &websocket.EventModel{
		ReceiverId: userId,
		Type:       websocket.EVT_SYNC,
		Payload:    nil,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(syncEvent)

	return nil
}

func (s *Service) UnrestrictUser(userId string) error {
	if err := s.userStore.ForceLiftRestriction(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     "Restriction Lifted",
		Content:   "The restriction to your account has been lifted by the administrator. Avoid committing violations to prevent future restrictions.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	}

	notificationHeading := "Account Restriction Lifted!"
	notificationContent := "Your account restriction has been lifted."

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(userId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: userId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	syncEvent := &websocket.EventModel{
		ReceiverId: userId,
		Type:       websocket.EVT_SYNC,
		Payload:    nil,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(syncEvent)

	return nil
}
