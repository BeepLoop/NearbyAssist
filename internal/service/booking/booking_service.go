package booking_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	notification_repo "nearbyassist/internal/repository/notification"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	qr_service "nearbyassist/internal/service/qr"
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
	ERR_FULLY_BOOKED             = "vendor is fully booked on the given schedule"
	ERR_INVALID_DATA             = "invalid data submitted"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	vendorStore  vendor_repo.VendorRepository
	notifStore   notification_repo.NotificationRepository
	bookingStore booking_repo.BookingRepository
	qrService    *qr_service.Service
	ws           websocket.Socket
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	serviceStore service_repo.ServiceRepository,
	vendorStore vendor_repo.VendorRepository,
	notifStore notification_repo.NotificationRepository,
	bookingStore booking_repo.BookingRepository,
	qrService *qr_service.Service,
	ws websocket.Socket,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		serviceStore: serviceStore,
		vendorStore:  vendorStore,
		notifStore:   notifStore,
		bookingStore: bookingStore,
		qrService:    qrService,
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

	if req.Quantity < 1 {
		return "", errors.New(ERR_INVALID_DATA)
	}

	if service.PricingType != models.FIXED_PRICING {
		cost := float64(req.Quantity) * utils.StringToFloat64ElseZero(service.Price)
		req.Cost = utils.Float64ToString(cost)
	}

	booking := &models.BookingModel{
		ClientId:           req.ClientId,
		VendorId:           req.VendorId,
		ServiceId:          req.ServiceId,
		ServiceTitle:       service.Title,
		ServiceDescription: service.Description,
		Price:              utils.Float64ToString(utils.StringToFloat64ElseZero(service.Price)),
		PricingType:        service.PricingType,
		Quantity:           req.Quantity,
		Cost:               req.Cost,
		Extras: slices.AppendSeq(
			make([]*models.BookingExtraModel, 0),
			utils.Map(req.Extras, func(extra request.Extra) *models.BookingExtraModel {
				return &models.BookingExtraModel{
					ExtraTitle:       utils.Must(s.encrypt.EncryptString(extra.Title)),
					ExtraDescription: utils.Must(s.encrypt.EncryptString(extra.Description)),
					Price:            extra.Price,
				}
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

	createdBooking, err := s.bookingStore.FindById(bookingId)
	if err != nil {
		return "", err
	}

	createdBookingPayload := response.Booking{
		Id: createdBooking.Id,
		Vendor: response.User{
			Id:       createdBooking.VendorId,
			Name:     utils.Must(s.encrypt.DecryptString(createdBooking.Vendor.Name)),
			ImageURL: createdBooking.Vendor.ImageUrl,
		},
		Client: response.User{
			Id:       createdBooking.ClientId,
			Name:     utils.Must(s.encrypt.DecryptString(createdBooking.Client.Name)),
			ImageURL: createdBooking.Client.ImageUrl,
		},
		ServiceId:          createdBooking.ServiceId,
		ServiceTitle:       utils.Must(s.encrypt.DecryptString(createdBooking.ServiceTitle)),
		ServiceDescription: utils.Must(s.encrypt.DecryptString(createdBooking.ServiceDescription)),
		Price:              createdBooking.Price,
		PricingType:        string(createdBooking.PricingType),
		Quantity:           createdBooking.Quantity,
		Cost:               createdBooking.Cost,
		Extras: slices.AppendSeq(
			make([]response.BookingExtra, 0),
			utils.Map(createdBooking.Extras, func(x *models.BookingExtraModel) response.BookingExtra {
				return response.BookingExtra{
					BookingId:   x.BookingId,
					Title:       x.ExtraTitle,
					Description: x.ExtraDescription,
					Price:       x.Price,
				}
			}),
		),
		Status:        string(createdBooking.Status),
		CreatedAt:     createdBooking.CreatedAt,
		UpdatedAt:     createdBooking.UpdatedAt,
		ScheduleStart: createdBooking.ScheduleStart.String,
		ScheduleEnd:   createdBooking.ScheduleEnd.String,
		CancelledBy:   createdBooking.CancelledBy.String,
		CancelReason:  createdBooking.CancelReason.String,
		QRSignature: utils.Must(s.qrService.SignData(&request.QRSignatureInput{
			ClientID:  createdBooking.ClientId,
			VendorID:  createdBooking.VendorId,
			BookingID: createdBooking.Id,
		})),
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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.VendorId,
		Type:       websocket.EVT_RECEIVED_BOOKING,
		Payload:    createdBookingPayload,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

	return bookingId, nil
}

func (s *Service) GetBooking(bookingId string) (*models.BookingModel, error) {
	booking, err := s.bookingStore.FindById(bookingId)
	if err != nil {
		return nil, err
	}

	booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
	booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
	booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
	booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

	if booking.Status == models.BOOKING_STATUS_CANCELLED || booking.Status == models.BOOKING_STATUS_REJECTED {
		booking.CancelReason.String = utils.Must(s.encrypt.DecryptString(booking.CancelReason.String))
		booking.CancelReason.Valid = true
	}

	for _, extra := range booking.Extras {
		extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
		extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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

	schedule, err := utils.StringToDateTime(booking.ScheduleStart.String)
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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_VENDOR_CANCELLED_BOOKING,
		Payload:    utils.Mapper{"id": req.BookingId, "reason": req.Reason},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.VendorId,
		Type:       websocket.EVT_CLIENT_CANCELLED_BOOKING,
		Payload:    utils.Mapper{"id": req.BookingId, "reason": req.Reason},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

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

	scheduleStart, err := utils.ParseDateString(req.ScheduleStart)
	if err != nil {
		return err
	}

	scheduleEnd, err := utils.ParseDateString(req.ScheduleEnd)
	if err != nil {
		return err
	}

	if err := utils.ValidateDateRange(scheduleStart, scheduleEnd); err != nil {
		return err
	}

	if available, err := s.vendorStore.IsDateAvailable(booking.VendorId, scheduleStart, scheduleEnd); err != nil {
		return err
	} else {
		if !available {
			return errors.New(ERR_FULLY_BOOKED)
		}
	}

	if err := s.bookingStore.Accept(req.BookingId, scheduleStart, scheduleEnd); err != nil {
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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_BOOKING_CONFIRMED,
		Payload: utils.Mapper{
			"id":            req.BookingId,
			"scheduleStart": utils.FormatDateTime(scheduleStart),
			"scheduleEnd":   utils.FormatDateTime(scheduleEnd),
		},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_BOOKING_REJECTED,
		Payload:    utils.Mapper{"id": req.BookingId, "reason": req.Reason},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

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
		booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
		booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
		booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
		booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

		for _, extra := range booking.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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
		booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
		booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
		booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
		booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

		for _, extra := range booking.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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
		booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
		booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
		booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
		booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

		for _, extra := range booking.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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
		booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
		booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
		booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
		booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

		for _, extra := range booking.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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
		reviewable.Vendor.Name = utils.Must(s.encrypt.DecryptString(reviewable.Vendor.Name))
		reviewable.Client.Name = utils.Must(s.encrypt.DecryptString(reviewable.Client.Name))
		reviewable.ServiceTitle = utils.Must(s.encrypt.DecryptString(reviewable.ServiceTitle))
		reviewable.ServiceDescription = utils.Must(s.encrypt.DecryptString(reviewable.ServiceDescription))

		for _, extra := range reviewable.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
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
		booking.Vendor.Name = utils.Must(s.encrypt.DecryptString(booking.Vendor.Name))
		booking.Client.Name = utils.Must(s.encrypt.DecryptString(booking.Client.Name))
		booking.ServiceTitle = utils.Must(s.encrypt.DecryptString(booking.ServiceTitle))
		booking.ServiceDescription = utils.Must(s.encrypt.DecryptString(booking.ServiceDescription))

		if booking.Status == models.BOOKING_STATUS_CANCELLED || booking.Status == models.BOOKING_STATUS_REJECTED {
			booking.CancelReason.String = utils.Must(s.encrypt.DecryptString(booking.CancelReason.String))
			booking.CancelReason.Valid = true
		}

		for _, extra := range booking.Extras {
			extra.ExtraTitle = utils.Must(s.encrypt.DecryptString(extra.ExtraTitle))
			extra.ExtraDescription = utils.Must(s.encrypt.DecryptString(extra.ExtraDescription))
		}
	}

	return bookings, nil
}

func (s *Service) CompleteBooking(bearerToken, bookingId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	booking, err := s.bookingStore.FindById(bookingId)
	if err != nil {
		return err
	} else {
		if booking.VendorId != userId {
			return errors.New(ERR_UNAUTHORIZED)
		}
	}

	if err := s.bookingStore.MarkComplete(bookingId); err != nil {
		return err
	}

	notificationHeading := "Booking complete"
	notificationContent := fmt.Sprintf(
		"Your booking with %s was completed",
		utils.Must(s.encrypt.DecryptString(booking.Vendor.Name)),
	)

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "success",
		Title:     notificationHeading,
		Content:   notificationContent,
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
	bookingEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_BOOKING_COMPLETE,
		Payload:    booking.Id,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(bookingEvent)

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

	scheduleStart, err := utils.ParseDateString(req.ScheduleStart)
	if err != nil {
		return err
	}

	scheduleEnd, err := utils.ParseDateString(req.ScheduleEnd)
	if err != nil {
		return err
	}

	if err := utils.ValidateDateRange(scheduleStart, scheduleEnd); err != nil {
		return err
	}

	if available, err := s.vendorStore.IsDateAvailable(booking.VendorId, scheduleStart, scheduleEnd); err != nil {
		return err
	} else {
		if !available {
			return errors.New(ERR_FULLY_BOOKED)
		}
	}

	if err := s.bookingStore.Reschedule(req.BookingId, scheduleStart, scheduleEnd); err != nil {
		return err
	}

	notificationHeading := "Booking rescheduled"
	notificationContent := "Your booking has been rescheduled"

	notification := &models.NotificationModel{
		Recipient: booking.ClientId,
		Type:      "generic",
		Title:     "Booking has been rescheduled",
		Content: fmt.Sprintf(
			"Your booking with the vendor: %s, has been rescheduled to %s - %s",
			utils.Must(s.encrypt.DecryptString(booking.Vendor.Name)),
			utils.FormatDateTime(scheduleStart),
			utils.FormatDateTime(scheduleEnd),
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
	reschedEvent := &websocket.EventModel{
		ReceiverId: booking.ClientId,
		Type:       websocket.EVT_BOOKING_RESCHEDULED,
		Payload: utils.Mapper{
			"id":            req.BookingId,
			"scheduleStart": utils.FormatDateTime(scheduleStart),
			"scheduleEnd":   utils.FormatDateTime(scheduleEnd),
		},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(reschedEvent)

	return nil
}
