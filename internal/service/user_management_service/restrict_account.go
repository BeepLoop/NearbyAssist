package user_management_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"time"
)

func (s *Service) RestrictUser(adminId, password, userId, reason, duration string) error {
	admin, err := s.adminStore.FindById(adminId)
	if err != nil {
		return err
	}

	if !core.IsPasswordMatch(admin.Password, password) {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if isRestricted, _, err := s.userStore.IsRestricted(userId); err != nil {
		return err
	} else {
		if isRestricted {
			return nil
		}
	}

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
		Title:     "Account Suspended " + stringifiedDuration,
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

	notificationHeading := "Account Suspended!"
	notificationContent := "You commited a violation resulting to account suspension."

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

func (s *Service) UnrestrictUser(adminId, password, userId string) error {
	admin, err := s.adminStore.FindById(adminId)
	if err != nil {
		return err
	}

	if !core.IsPasswordMatch(admin.Password, password) {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if isRestricted, _, err := s.userStore.IsRestricted(userId); err != nil {
		return err
	} else {
		if !isRestricted {
			return nil
		}
	}

	if err := s.userStore.ForceLiftRestriction(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     "Suspension Lifted",
		Content:   "Suspension to your account has been lifted. Avoid violations of community guidelines to prevent future suspensions.",
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

	notificationHeading := "Account Suspension Lifted!"
	notificationContent := "Your account suspension has been lifted."

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
