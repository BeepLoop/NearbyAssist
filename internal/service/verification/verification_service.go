package verification_service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	user_repo "nearbyassist/internal/repository/user"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
)

type Service struct {
	userStore         user_repo.UserRepository
	verificationStore verification_repo.VerificationRepository
	notifStore        notification_repo.NotificationRepository
	ws                websocket.Socket
	fs                fs.FileStorage
	encrypt           core.Encryption
	jwt               core.Authenticator
}

func NewService(userStore user_repo.UserRepository, verificationStore verification_repo.VerificationRepository, notifStore notification_repo.NotificationRepository, ws websocket.Socket, fs fs.FileStorage, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{
		userStore:         userStore,
		verificationStore: verificationStore,
		notifStore:        notifStore,
		ws:                ws,
		fs:                fs,
		encrypt:           encrypt,
		jwt:               jwt,
	}
}

func (s *Service) CreateVerificationRequest(name, phone, address, idType, idNumber, bearerToken string, latitude, longitude float64, files []*multipart.FileHeader) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	encryptedName, err := s.encrypt.EncryptString(name)
	if err != nil {
		return "", err
	}

	encryptedPhone, err := s.encrypt.EncryptString(phone)
	if err != nil {
		return "", err
	}

	encryptedAddress, err := s.encrypt.EncryptString(address)
	if err != nil {
		return "", err
	}

	encryptedIdNumber, err := s.encrypt.EncryptString(idNumber)
	if err != nil {
		return "", err
	}

	req := new(models.IdentityVerificationModel)
	req.UserId = userId
	req.Name = encryptedName
	req.Address = encryptedAddress
	req.Phone = encryptedPhone
	req.IdType = idType
	req.IdNumber = encryptedIdNumber
	req.Latitude = latitude
	req.Longitude = longitude

	for _, file := range files {
		// Read bytes
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return "", err
		}

		// Encrypt the file
		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		switch file.Filename {
		case "frontId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_FRONT,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.FrontIdImageUrl = url
			}

		case "backId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_BACK,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.BackIdImageUrl = url
			}

		case "face":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.FACE,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.FaceImageUrl = url
			}

		default:
			return "", err
		}
	}

	verificationId, err := s.verificationStore.Create(req)
	if err != nil {
		return "", err
	}

	return verificationId, nil
}

func (s *Service) GetRequest(id string) (*models.IdentityVerificationModel, error) {
	request, err := s.verificationStore.FindById(id)
	if err != nil {
		return nil, err
	}

	decryptedName, err := s.encrypt.DecryptString(request.Name)
	if err != nil {
		return nil, err
	}
	request.Name = decryptedName

	decryptedAddress, err := s.encrypt.DecryptString(request.Address)
	if err != nil {
		return nil, err
	}
	request.Address = decryptedAddress

	decryptedIdNumber, err := s.encrypt.DecryptString(request.IdNumber)
	if err != nil {
		return nil, err
	}
	request.IdNumber = decryptedIdNumber

	return request, nil
}

func (s *Service) GetIdentityVerificationRequests() ([]*models.IdentityVerificationModel, error) {
	requests, err := s.verificationStore.GetAll("pending")
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		decryptedName, err := s.encrypt.DecryptString(request.Name)
		if err != nil {
			return nil, err
		}
		request.Name = decryptedName

		decryptedAddress, err := s.encrypt.DecryptString(request.Address)
		if err != nil {
			return nil, err
		}
		request.Address = decryptedAddress

		decryptedIdNumber, err := s.encrypt.DecryptString(request.IdNumber)
		if err != nil {
			return nil, err
		}
		request.IdNumber = decryptedIdNumber
	}

	return requests, nil
}

func (s *Service) AcceptRequest(id string) error {
	request, err := s.verificationStore.FindById(id)
	if err != nil {
		return err
	}

	if err := s.verificationStore.AcceptRequest(id); err != nil {
		return err
	}

	notificationHeading := "Identity Verification Accepted"
	notificationContent := "Congratulations! Your identity verification request is approved."

	notification := &models.NotificationModel{
		Recipient: request.UserId,
		Type:      "success",
		Title:     "Identity Verification Accepted",
		Content:   "Congratulations! Your identity verification request has been accepted. Go to your settings and Sync Account to see the changes.",
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

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(request.UserId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: request.UserId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	// send sync event to instruct client to pull the udpated values
	syncEvent := &websocket.EventModel{
		ReceiverId: request.UserId,
		Type:       websocket.EVT_SYNC,
		Payload:    nil,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(syncEvent)

	return nil
}

func (s *Service) RejectRequest(id, reason string) error {
	if reason == "" {
		return errors.New("invalid reason")
	}

	encryptedReason, err := s.encrypt.EncryptString(reason)
	if err != nil {
		return err
	}

	if err := s.verificationStore.RejectRequest(id, encryptedReason); err != nil {
		return err
	}

	request, err := s.verificationStore.FindById(id)
	if err != nil {
		return err
	}

	notificationHeading := "Identity Verification Rejected"
	notificationContent := "Your verification request is rejected"

	notification := &models.NotificationModel{
		Recipient: request.UserId,
		Type:      "fail",
		Title:     notificationHeading,
		Content:   "Your identity verification request is rejected. Reason of rejection: " + reason,
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

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(request.UserId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	event := &websocket.EventModel{
		ReceiverId: request.UserId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(event)

	return nil
}
