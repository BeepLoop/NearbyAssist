package booking_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	notification_repo "nearbyassist/internal/repository/notification"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	ERR_DISABLED_SERVICE         = "service is disabled"
	ERR_HAS_PENDING_OR_CONFIRMED = "already have pending or confirmed booking for service"
	ERR_DISALLOWED_ACTION        = "invalid action performed"
	ERR_UNAUTHORIZED             = "unauthorized"
	ERR_SCHEDULE_OVERLAP         = "schedule overlap"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	notifStore   notification_repo.NotificationRepository
	bookingStore booking_repo.BookingRepository
	ws           websocket.Socket
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	serviceStore service_repo.ServiceRepository,
	notifStore notification_repo.NotificationRepository,
	bookingStore booking_repo.BookingRepository,
	ws websocket.Socket,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		serviceStore: serviceStore,
		notifStore:   notifStore,
		bookingStore: bookingStore,
		ws:           ws,
		encrypt:      encrypt,
		jwt:          jwt,
	}
}

func (s *Service) CreateBooking(req *request.NewBookingPayload) (string, error) {
	service, err := s.serviceStore.FindById(req.ServiceId)
	if err != nil {
		return "", err
	}
	if service.Disabled {
		return "", errors.New(ERR_DISABLED_SERVICE)
	}

	booking := &models.BookingModel{
		ClientId:  req.ClientId,
		VendorId:  req.VendorId,
		ServiceId: req.ServiceId,
		Cost:      req.Cost,
		Extras: slices.AppendSeq(
			make([]*models.ExtraModel, 0),
			utils.Map(req.Extras, func(extra request.Extra) *models.ExtraModel {
				return &models.ExtraModel{Model: models.Model{Id: extra.Id}}
			}),
		),
	}

	hasOngoing, err := s.bookingStore.HasOngoingBookingForService(booking)
	if err != nil {
		return "", err
	}
	if hasOngoing {
		return "", errors.New(ERR_HAS_PENDING_OR_CONFIRMED)
	}

	bookingId, err := s.bookingStore.Create(booking)
	if err != nil {
		return "", err
	}

	notificationHeading := "New Request"
	notificationContent := "1 new booking request"

	notification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "generic",
		Title:     "New Request",
		Content:   "You received a booking request. View reqeust in your booking dashboard.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "generic",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return "", err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return bookingId, nil
}

func (s *Service) GetBooking(bookingId string) (*models.BookingModel, error) {
	booking, err := s.bookingStore.FindById(bookingId)
	if err != nil {
		return nil, err
	}

	booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
	booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
	booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
	booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

	if booking.Status == models.BOOKING_STATUS_CANCELLED {
		booking.CancelReason.String = utils.Must(s.encrypt.DecryptString(booking.CancelReason.String))
		booking.CancelReason.Valid = true
	}

	for _, extra := range booking.Extras {
		extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
		extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
	}

	return booking, nil
}

func (s *Service) VendorCancelBooking(bearerToken string, req *request.CancelBookingPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}
	if booking.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}
	if booking.Status != models.BOOKING_STATUS_CONFIRMED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}
	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	schedule, err := utils.StringToDateTime(booking.ScheduledAt.String)
	if err != nil {
		return err
	}
	now, err := utils.StringToDateTime(utils.CurrentTimeStamp())
	if err != nil {
		return err
	}
	if !now.After(schedule) {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))
	if err := s.bookingStore.Cancel(req.BookingId, userId, encryptedReason); err != nil {
		return err
	}

	notificationHeading := "Booking Cancelled"
	notificationContent := "Vendor cancelled your booking with them"

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "fail",
		Title:     "Scheduled booking cancelled",
		Content:   fmt.Sprintf("The vendor cancelled your scheduled booking with them. Rason: %s", req.Reason),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) ClientCancelBooking(bearerToken string, req *request.CancelBookingPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}

	if booking.Status != models.BOOKING_STATUS_PENDING {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.ClientId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))
	if err := s.bookingStore.Cancel(req.BookingId, userId, encryptedReason); err != nil {
		return err
	}

	notificationHeading := "Booking Request Cancelled"
	notificationContent := "A client cancelled their booking request"

	notification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "fail",
		Title:     "Booking Request Cancelled",
		Content:   fmt.Sprintf("A client cancelled their request for your service. Rason: %s", req.Reason),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) AcceptBookingRequest(bearerToken string, req *request.AcceptBookingPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}

	if booking.Status != models.BOOKING_STATUS_PENDING {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	schedule := utils.FormatDate(req.Schedule)
	if err := s.bookingStore.Accept(req.BookingId, schedule); err != nil {
		return err
	}

	notificationHeading := "Booking Request Accepted"
	notificationContent := "Your booking request was accepted by the vendor"

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "success",
		Title:     "Booking Request Accepted",
		Content:   "Your booking request has been accepted by the vendor",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) RejectBookingRequest(bearerToken string, req *request.RejectRequestPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}

	if booking.Status != models.BOOKING_STATUS_PENDING {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	if booking.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))
	if err := s.bookingStore.Reject(req.BookingId, encryptedReason); err != nil {
		return err
	}

	notificationHeading := "Booking Request Rejected"
	notificationContent := "Your booking request was rejected by the vendor"

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "fail",
		Title:     "Booking Request Rejected",
		Content:   fmt.Sprintf("Your booking request was rejected by the vendor. Reason: %s", req.Reason),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) GetBookingUserSent(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetBookingSent(userId)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
		booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
		booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
		booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

		for _, extra := range booking.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return bookings, nil
}

func (s *Service) GetBookingUserReceived(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetBookingReceived(userId)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
		booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
		booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
		booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

		for _, extra := range booking.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return bookings, nil
}

func (s *Service) GetRecentBookings(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetRecent(userId)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
		booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
		booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
		booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

		for _, extra := range booking.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return bookings, nil
}

func (s *Service) GetConfirmedBookings(bearerToken, filter string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetConfirmed(userId, filter)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
		booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
		booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
		booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

		for _, extra := range booking.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return bookings, nil
}

func (s *Service) GetReviewableBookings(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	reviewables, err := s.bookingStore.GetReviewableBookings(userId)
	if err != nil {
		return nil, err
	}

	for _, reviewable := range reviewables {
		reviewable.Vendor = utils.Must(s.encrypt.DecryptString(reviewable.Vendor))
		reviewable.Client = utils.Must(s.encrypt.DecryptString(reviewable.Client))
		reviewable.Service.Title = utils.Must(s.encrypt.DecryptString(reviewable.Service.Title))
		reviewable.Service.Description = utils.Must(s.encrypt.DecryptString(reviewable.Service.Description))

		for _, extra := range reviewable.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return reviewables, nil
}

func (s *Service) GetBookingHistory(bearerToken, filter string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetHistory(userId, filter)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		booking.Vendor = utils.Must(s.encrypt.DecryptString(booking.Vendor))
		booking.Client = utils.Must(s.encrypt.DecryptString(booking.Client))
		booking.Service.Title = utils.Must(s.encrypt.DecryptString(booking.Service.Title))
		booking.Service.Description = utils.Must(s.encrypt.DecryptString(booking.Service.Description))

		if booking.Status == models.BOOKING_STATUS_CANCELLED {
			booking.CancelReason.String = utils.Must(s.encrypt.DecryptString(booking.CancelReason.String))
			booking.CancelReason.Valid = true
		}

		for _, extra := range booking.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return bookings, nil
}

func (s *Service) CompleteBooking(bearerToken, bookingId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if booking, err := s.bookingStore.FindById(bookingId); err != nil {
		return err
	} else {
		if booking.VendorId != userId {
			return errors.New(ERR_UNAUTHORIZED)
		}
	}

	if err := s.bookingStore.MarkComplete(bookingId); err != nil {
		return err
	}

	return nil
}

func (s *Service) Reschedule(bearerToken string, req *request.RescheduleBookingPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}

	if booking.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New(ERR_DISALLOWED_ACTION)
	}

	confirmedBookings, err := s.bookingStore.GetConfirmedBookingsOfVendor(userId)
	if err != nil {
		return err
	}
	if utils.HasScheduleOverlap(req.Schedule, confirmedBookings) {
		return errors.New(ERR_SCHEDULE_OVERLAP)
	}

	schedule := utils.FormatDate(req.Schedule)
	if err := s.bookingStore.Reschedule(req.BookingId, schedule); err != nil {
		return err
	}

	notificationHeading := "Booking rescheduled"
	notificationContent := "Your booking has been rescheduled"

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "generic",
		Title:     "Booking has been rescheduled",
		Content: fmt.Sprintf(
			"Your booking with the vendor: %s, has been rescheduled to %s",
			utils.Must(s.encrypt.DecryptString(booking.Vendor)),
			schedule,
		),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "generic",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(booking.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}
