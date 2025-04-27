package user_management_service

import (
	"fmt"
	"nearbyassist/internal/models"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"time"
)

func (s *Service) RestrictUser(userId, reason, duration string) error {
	d, err := utils.StringDaysToDuration(duration)
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

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
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
		Content:   "Restriction to your account has been lifted. Avoid violations of community guidelines to prevent future restrictions.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
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
