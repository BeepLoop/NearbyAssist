package application_service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/models"
	application_repo "nearbyassist/internal/repository/application"
	notification_repo "nearbyassist/internal/repository/notification"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/utils"
)

type Service struct {
	store      application_repo.ApplicationRepository
	notifStore notification_repo.NotificationRepository
	fs         fs.FileStorage
	encrypt    core.Encryption
	jwt        core.Authenticator
}

func NewService(store application_repo.ApplicationRepository, notifStore notification_repo.NotificationRepository, fs fs.FileStorage, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{store: store, notifStore: notifStore, fs: fs, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateApplication(bearerToken, expertiseId string, files []*multipart.FileHeader) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	application := new(models.ApplicationModel)
	application.ApplicantId = userId
	application.ExpertiseId = expertiseId

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return "", err
		}

		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		switch file.Filename {
		case "policeClearance":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.POLICE_CLEARANCE_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return "", err
			}

			application.PoliceClearanceUrl = url

		case "supportingDocument":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.APPLICATION_PROOF_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return "", err
			}

			application.SupportingDocumentUrl = url

		default:
			return "", err

		}
	}

	applicationId, err := s.store.Create(application)
	if err != nil {
		return "", err
	}

	return applicationId, nil
}

func (s *Service) GetApplications() ([]*models.ApplicationModel, error) {
	applications, err := s.store.GetAll("pending")
	if err != nil {
		return nil, err
	}

	for _, application := range applications {
		if decrypted, err := s.encrypt.DecryptString(application.ApplicantName); err != nil {
			return nil, err
		} else {
			application.ApplicantName = decrypted
		}
	}

	return applications, nil
}

func (s *Service) GetApplicationDetail(applicationId string) (*models.ApplicationModel, error) {
	application, err := s.store.FindById(applicationId)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(application.ApplicantName); err != nil {
		return nil, err
	} else {
		application.ApplicantName = decrypted
	}

	return application, nil
}

func (s *Service) AcceptRequest(applicationId string) error {
	application, err := s.store.FindById(applicationId)
	if err != nil {
		return err
	}

	if decrypted, err := s.encrypt.DecryptString(application.ApplicantName); err != nil {
		return err
	} else {
		application.ApplicantName = decrypted
	}

	if err := s.store.AcceptRequest(applicationId); err != nil {
		return err
	}

	notificationHeading := "Expertise request granted"
	notificationContent := "Congratulations! Your expertise request has been granted."

	notification := &models.NotificationModel{
		Recipient: application.ApplicantId,
		Type:      "success",
		Title:     "Sucessfully added an expertise",
		Content:   "Congratulations! You successfully added expertise in: " + application.Expertise,
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
		if err := oneSignal.NewUrgentNotification(application.ApplicantId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

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

	if err := s.store.RejectRequest(id, encryptedReason); err != nil {
		return err
	}

	application, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	notificationHeading := "Expertise request denied"
	notificationContent := "We are sorry to inform that your request is rejected."

	notification := &models.NotificationModel{
		Recipient: application.ApplicantId,
		Type:      "fail",
		Title:     notificationHeading,
		Content:   fmt.Sprintf("We are sorry to inform you that your request to add expertise in %s is denied. Reason: %s", application.Expertise, reason),
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
		if err := oneSignal.NewUrgentNotification(application.ApplicantId, notificationHeading, notificationContent); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	return nil
}
