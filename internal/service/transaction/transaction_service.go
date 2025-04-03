package transaction_service

import (
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	transaction_repo "nearbyassist/internal/repository/transaction"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
)

type Service struct {
	notifStore       notification_repo.NotificationRepository
	transactionStore transaction_repo.TransactionRepository
	ws               websocket.Socket
	encrypt          core.Encryption
	jwt              core.Authenticator
}

func NewService(
	notifStore notification_repo.NotificationRepository,
	transactionStore transaction_repo.TransactionRepository,
	ws websocket.Socket,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		notifStore:       notifStore,
		transactionStore: transactionStore,
		ws:               ws,
		encrypt:          encrypt,
		jwt:              jwt,
	}
}

func (s *Service) CreateTransaction(req *request.NewTransactionPayload) (string, error) {
	transaction := &models.TransactionModel{
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
	transaction.Extras = extras

	transactionId, err := s.transactionStore.Create(transaction)
	if err != nil {
		return "", err
	}

	notificationHeading := "New Request"
	notificationContent := "1 new transaction request"

	notification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "generic",
		Title:     "New Request",
		Content:   "You received a transaction request. View reqeust in your transaction dashboard.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "generic",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return "", err
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(transaction.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: transaction.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return transactionId, nil
}

func (s *Service) GetTransaction(transactionId string) (*models.TransactionModel, error) {
	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return nil, err
	}

	transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
	transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
	transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
	transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

	if transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		transaction.CancelReason.String = utils.Must(s.encrypt.DecryptString(transaction.CancelReason.String))
		transaction.CancelReason.Valid = true
	}

	for _, extra := range transaction.Extras {
		extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
		extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
	}

	return transaction, nil
}

func (s *Service) CancelTransaction(bearerToken string, req *request.CancelRequestPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(req.TransactionId)
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

	encryptedReason := utils.Must(s.encrypt.EncryptString(req.Reason))

	if err := s.transactionStore.Cancel(req.TransactionId, encryptedReason); err != nil {
		return err
	}

	notificationHeading := "Transaction Request Cancelled"
	notificationContent := "A client cancelled their transaction request"

	notification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "fail",
		Title:     "Transaction Request Cancelled",
		Content:   fmt.Sprintf("A client cancelled their request for your service. Rason: %s", req.Reason),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(transaction.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: transaction.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) AcceptTransactionRequest(bearerToken string, req *request.AcceptTransactionPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(req.TransactionId)
	if err != nil {
		return err
	}

	if transaction.Status != models.TRANSACTION_STATUS_PENDING {
		return errors.New("Could not accept non-pending transaction")
	}

	if transaction.Status == models.TRANSACTION_STATUS_DONE || transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		return errors.New("Transaction already completed or cancelled")
	}

	if transaction.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	schedule := utils.FormatDate(req.Schedule)
	if err := s.transactionStore.Accept(req.TransactionId, schedule); err != nil {
		return err
	}

	notificationHeading := "Transaction Request Accepted"
	notificationContent := "Your transaction request was accepted by the vendor"

	notification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "success",
		Title:     "Transaction Request Accepted",
		Content:   "Your transaction request has been accepted by the vendor",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(transaction.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: transaction.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) RejectTransactionRequest(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	transaction, err := s.transactionStore.FindById(transactionId)
	if err != nil {
		return err
	}

	if transaction.Status != models.TRANSACTION_STATUS_PENDING {
		return errors.New("Could not reject non-pending transaction")
	}

	if transaction.Status == models.TRANSACTION_STATUS_DONE || transaction.Status == models.TRANSACTION_STATUS_CANCELLED {
		return errors.New("Transaction already completed or cancelled")
	}

	if transaction.VendorId != userId {
		return errors.New("Unauthorized accept request")
	}

	if err := s.transactionStore.Reject(transactionId); err != nil {
		return err
	}

	notificationHeading := "Transaction Request Rejected"
	notificationContent := "Your transaction request was rejected by the vendor"

	notification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "fail",
		Title:     "Transaction Request Rejected",
		Content:   "Your transaction request was rejected by the vendor",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: transaction.VendorId,
		Type:      "fail",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(transaction.ClientId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: transaction.ClientId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(notifEvent)

	return nil
}

func (s *Service) GetTransactionUserSent(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetTransactionSent(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
		transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
		transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
		transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

		for _, extra := range transaction.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return transactions, nil
}

func (s *Service) GetTransactionUserReceived(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetTransactionReceived(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
		transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
		transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
		transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

		for _, extra := range transaction.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return transactions, nil
}

func (s *Service) GetRecentTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetRecent(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
		transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
		transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
		transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

		for _, extra := range transaction.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return transactions, nil
}

func (s *Service) GetConfirmedTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetConfirmed(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
		transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
		transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
		transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

		for _, extra := range transaction.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return transactions, nil
}

func (s *Service) GetReviewableTransactions(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	reviewables, err := s.transactionStore.GetReviewableTransactions(userId)
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

func (s *Service) GetTransactionHistory(bearerToken string) ([]*models.TransactionModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionStore.GetHistory(userId)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Vendor = utils.Must(s.encrypt.DecryptString(transaction.Vendor))
		transaction.Client = utils.Must(s.encrypt.DecryptString(transaction.Client))
		transaction.Service.Title = utils.Must(s.encrypt.DecryptString(transaction.Service.Title))
		transaction.Service.Description = utils.Must(s.encrypt.DecryptString(transaction.Service.Description))

		for _, extra := range transaction.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return transactions, nil
}

func (s *Service) CompleteTransaction(bearerToken, transactionId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if transaction, err := s.transactionStore.FindById(transactionId); err != nil {
		return err
	} else {
		if transaction.VendorId != userId {
			return errors.New("unauthorized")
		}
	}

	if err := s.transactionStore.MarkComplete(transactionId); err != nil {
		return err
	}

	return nil
}
