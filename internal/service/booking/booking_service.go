package booking_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	notification_repo "nearbyassist/internal/repository/notification"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
)

type Service struct {
	notifStore   notification_repo.NotificationRepository
	bookingStore booking_repo.BookingRepository
	ws           websocket.Socket
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	notifStore notification_repo.NotificationRepository,
	bookingStore booking_repo.BookingRepository,
	ws websocket.Socket,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		notifStore:   notifStore,
		bookingStore: bookingStore,
		ws:           ws,
		encrypt:      encrypt,
		jwt:          jwt,
	}
}

func (s *Service) CreateBooking(req *request.NewBookingPayload) (string, error) {
	booking := &models.BookingModel{
		ClientId:  req.ClientId,
		VendorId:  req.VendorId,
		ServiceId: req.ServiceId,
		Cost:      req.Cost,
	}

	extras := make([]*models.ExtraModel, 0)
	for _, extra := range req.Extras {
		extras = append(extras, &models.ExtraModel{
			Model: models.Model{
				Id: extra.Id,
			},
		})
	}
	booking.Extras = extras

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

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return "", err
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

func (s *Service) CancelBooking(bearerToken string, req *request.CancelRequestPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return err
	}

	if booking.Status != models.BOOKING_STATUS_PENDING {
		return errors.New("Could not cancel non-pending booking")
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New("Booking already completed or cancelled")
	}

	if booking.ClientId != userId {
		return errors.New("Unauthorized cancel request")
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))

	if err := s.bookingStore.Cancel(req.BookingId, encryptedReason); err != nil {
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

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
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
		return errors.New("Could not accept non-pending booking")
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New("Booking already completed or cancelled")
	}

	if booking.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	schedule := utils.FormatDate(req.Schedule)
	if err := s.bookingStore.Accept(req.BookingId, schedule); err != nil {
		return err
	}

	notificationHeading := "Booking Request Accepted"
	notificationContent := "Your booking request was accepted by the vendor"

	notification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "success",
		Title:     "Booking Request Accepted",
		Content:   "Your booking request has been accepted by the vendor",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
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
		return errors.New("Could not reject non-pending booking")
	}

	if booking.Status == models.BOOKING_STATUS_DONE || booking.Status == models.BOOKING_STATUS_CANCELLED {
		return errors.New("Booking already completed or cancelled")
	}

	if booking.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))
	if err := s.bookingStore.Reject(req.BookingId, encryptedReason); err != nil {
		return err
	}

	notificationHeading := "Booking Request Rejected"
	notificationContent := "Your booking request was rejected by the vendor"

	notification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "fail",
		Title:     "Booking Request Rejected",
		Content:   fmt.Sprintf("Your booking request was rejected by the vendor. Reason: %s", req.Reason),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: booking.VendorId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
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

func (s *Service) GetConfirmedBookings(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetConfirmed(userId)
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

func (s *Service) GetBookingHistory(bearerToken string) ([]*models.BookingModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetHistory(userId)
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

func (s *Service) CompleteBooking(bearerToken, bookingId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if booking, err := s.bookingStore.FindById(bookingId); err != nil {
		return err
	} else {
		if booking.VendorId != userId {
			return errors.New("unauthorized")
		}
	}

	if err := s.bookingStore.MarkComplete(bookingId); err != nil {
		return err
	}

	return nil
}
