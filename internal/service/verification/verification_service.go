package verification_service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	user_repo "nearbyassist/internal/repository/user"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	notification_service "nearbyassist/internal/service/notification"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	ERR_ALREADY_VERIFIED = "account already verified"
	ERR_INVALID_REASON   = "provided reason not allowed"
)

type Service struct {
	userStore         user_repo.UserRepository
	verificationStore verification_repo.VerificationRepository
	notifStore        notification_repo.NotificationRepository
	resourceService   *resource_service.Service
	ws                websocket.Socket
	fs                fs.FileStorage
	encrypt           core.Encryption
	jwt               core.Authenticator
}

func NewService(
	userStore user_repo.UserRepository,
	verificationStore verification_repo.VerificationRepository,
	notifStore notification_repo.NotificationRepository,
	resourceService *resource_service.Service,
	ws websocket.Socket,
	fs fs.FileStorage,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:         userStore,
		verificationStore: verificationStore,
		notifStore:        notifStore,
		resourceService:   resourceService,
		ws:                ws,
		fs:                fs,
		encrypt:           encrypt,
		jwt:               jwt,
	}
}

func (s *Service) VerifyAccount(bearerToken string, payload *request.VerifyAccountPayload, files []*multipart.FileHeader) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	user, err := s.userStore.FindById(userId)
	if err != nil {
		return err
	}
	if user.Verified {
		return errors.New(ERR_ALREADY_VERIFIED)
	}

	verificationRequest := &models.IdentityVerificationModel{
		UserId: user.Id,
		User: models.UserModel{
			Model: user.Model,
			Name:  utils.Must(s.encrypt.EncryptString(payload.Name)),
			Phone: utils.Must(s.encrypt.EncryptString(payload.Phone)),
			Address: models.AddressModel{
				Id:        user.Address.Id,
				Address:   utils.Must(s.encrypt.EncryptString(payload.Address)),
				Latitude:  payload.Latitude,
				Longitude: payload.Longitude,
			},
			Identification: models.IdentificationModel{
				Type:            payload.IdType,
				ReferenceNumber: utils.Must(s.encrypt.EncryptString(payload.ReferenceNumber)),
			},
		},
	}

	for _, file := range files {
		// Read bytes
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return err
		}

		// Encrypt the file
		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return err
		}

		switch file.Filename {
		case "frontId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_FRONT,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return err
			} else {
				verificationRequest.User.Identification.FrontImageUrl = url
			}

		case "backId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_BACK,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return err
			} else {
				verificationRequest.User.Identification.BackImageUrl = url
			}

		case "face":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.FACE,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return err
			} else {
				verificationRequest.User.Identification.SelfieImageUrl = url
			}

		default:
			return err
		}
	}

	if err := s.verificationStore.UpdateUserInfo(verificationRequest, false); err != nil {
		return err
	}

	if _, err := s.verificationStore.CreateLink(verificationRequest.UserId); err != nil {
		return err
	}

	sse.New().IncreaseVerification()

	return nil
}

func (s *Service) GetUserWithRequest(requestId string) (*models.UserModel, error) {
	request, err := s.verificationStore.FindById(requestId)
	if err != nil {
		return nil, err
	}

	user := &models.UserModel{
		Model: models.Model{Id: request.UserId},
	}

	return user, nil
}

func (s *Service) GetRequestList() ([]dto.VerificationRequest, error) {
	requests, err := s.verificationStore.GetAll(models.IDENTITY_VERIF_STATUS_PENDING)
	if err != nil {
		return nil, err
	}

	data := slices.AppendSeq(
		make([]dto.VerificationRequest, 0),
		utils.Map(requests, func(request *models.IdentityVerificationModel) dto.VerificationRequest {
			return dto.VerificationRequest{
				Id:              request.Id,
				UserID:          request.User.Id,
				Name:            utils.Must(s.encrypt.DecryptString(request.User.Name)),
				Email:           utils.Must(s.encrypt.DecryptString(request.User.Email)),
				ImageURL:        request.User.ImageUrl,
				Phone:           utils.Must(s.encrypt.DecryptString(request.User.Phone)),
				Address:         utils.Must(s.encrypt.DecryptString(request.User.Address.Address)),
				IDType:          request.User.Identification.Type,
				ReferenceNumber: request.User.Identification.ReferenceNumber,
				IDFrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.FrontImageUrl)),
				IDBackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.BackImageUrl)),
				SelfieImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.SelfieImageUrl)),
				CreatedAt:       utils.FormatDate(request.CreatedAt),
			}
		}),
	)

	return data, nil
}

func (s *Service) GetRequest(id string) (*dto.VerificationRequest, error) {
	request, err := s.verificationStore.FindById(id)
	if err != nil {
		return nil, err
	}

	data := &dto.VerificationRequest{
		Id:              request.Id,
		UserID:          request.User.Id,
		Name:            utils.Must(s.encrypt.DecryptString(request.User.Name)),
		Email:           utils.Must(s.encrypt.DecryptString(request.User.Email)),
		ImageURL:        request.User.ImageUrl,
		Phone:           utils.Must(s.encrypt.DecryptString(request.User.Phone)),
		Address:         utils.Must(s.encrypt.DecryptString(request.User.Address.Address)),
		IDType:          request.User.Identification.Type,
		ReferenceNumber: utils.Must(s.encrypt.DecryptString(request.User.Identification.ReferenceNumber)),
		IDFrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.FrontImageUrl)),
		IDBackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.BackImageUrl)),
		SelfieImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(request.User.Identification.SelfieImageUrl)),
		CreatedAt:       request.CreatedAt,
	}

	return data, nil
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
		Recipient: request.User.Id,
		Type:      "success",
		Title:     "Identity Verification Accepted",
		Content:   "Congratulations! Your identity verification request has been accepted. Go to your settings and Sync Account to see the changes.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: request.User.Id,
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
	if err := oneSignal.NewUrgentNotification(request.User.Id, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: request.User.Id,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	syncEvent := &websocket.EventModel{
		ReceiverId: request.User.Id,
		Type:       websocket.EVT_SYNC,
		Payload:    nil,
	}

	s.ws.Send(notifEvent)
	s.ws.Send(syncEvent)

	return nil
}

func (s *Service) RejectRequest(id, reason string) error {
	if reason == "" {
		return errors.New(ERR_INVALID_REASON)
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
		Recipient: request.User.Id,
		Type:      "fail",
		Title:     notificationHeading,
		Content:   "Your identity verification request is rejected. Reason of rejection: " + reason,
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: request.User.Id,
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
	if err := oneSignal.NewUrgentNotification(request.User.Id, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	event := &websocket.EventModel{
		ReceiverId: request.User.Id,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(event)

	return nil
}
