package verification_service

import (
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/utils"
)

type Service struct {
	store      verification_repo.VerificationRepository
	notifStore notification_repo.NotificationRepository
	fs         fs.FileStorage
	encrypt    core.Encryption
	jwt        core.Authenticator
}

func NewService(store verification_repo.VerificationRepository, notifStore notification_repo.NotificationRepository, fs fs.FileStorage, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{store: store, notifStore: notifStore, fs: fs, encrypt: encrypt, jwt: jwt}
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

	verificationId, err := s.store.Create(req)
	if err != nil {
		return "", err
	}

	return verificationId, nil
}

func (s *Service) GetRequest(id string) (*models.IdentityVerificationModel, error) {
	request, err := s.store.FindById(id)
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
	requests, err := s.store.GetAll("pending")
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
	request, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	if err := s.store.AcceptRequest(id); err != nil {
		return err
	}

	notificationHeading := "Identity Verification Accepted"
	notificationContent := "Congratulations! Your identity verification request has been accepted. Go to your settings and Sync Account to see the changes."

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

	oneSignal := notification_service.OneSignalInstance
	if oneSignal != nil {
		if err := oneSignal.NewUrgentNotification(request.UserId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	return nil
}

func (s *Service) RejectRequest(id, reason string) error {
	request, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	if err := s.store.RejectRequest(id); err != nil {
		return err
	}

	notificationHeading := "Identity Verification Rejected"
	notificationContent := "Identity Verification Rejected" + reason

	notification := &models.NotificationModel{
		Recipient: request.UserId,
		Type:      "fail",
		Title:     notificationHeading,
		Content:   notificationContent,
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

	oneSignal := notification_service.OneSignalInstance
	if oneSignal != nil {
		if err := oneSignal.NewUrgentNotification(request.UserId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	return nil
}
