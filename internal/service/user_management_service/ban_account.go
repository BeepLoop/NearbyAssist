package user_management_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	MINIMUM_REPORTS_FOR_BANNING = 5

	ERR_DID_NOT_MEET_BANNING_REQUIREMENT = "account did not meet minimum number of reports to ban"
)

func (s *Service) BanUser(userId string) error {
	reports, err := s.reportUserStore.GetAllReportedIs(userId)
	if err != nil {
		return err
	}

	resolvedReports := slices.Collect(utils.Retain(
		reports,
		func(report *models.UserReportModel) bool {
			return report.Status == models.REPORT_STATUS_RESOLVED
		}),
	)
	if len(resolvedReports) < MINIMUM_REPORTS_FOR_BANNING {
		return errors.New(ERR_DID_NOT_MEET_BANNING_REQUIREMENT)
	}

	if err := s.userStore.BanUser(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     "Your account has been suspended",
		Content:   "Your account has been indifinetly suspended due to multiple, repeated violations of community guidelines.",
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

	notificationHeading := "Account Banned"
	notificationContent := "Your account has been indifinetely suspended due to violations."

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

func (s *Service) UnbanUser(userId string) error {
	if err := s.userStore.UnbanUser(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     "Account suspension lifted",
		Content:   "You now have full access to your account and our services again. Make sure to follow community guidelines to avoid future suspensions. Welcome back!",
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

	notificationHeading := "Account suspension lifted"
	notificationContent := "Indifinite suspension of account has been lifted. Welcome back!"

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
