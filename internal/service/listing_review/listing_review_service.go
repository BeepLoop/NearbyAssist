package listingreview

import (
	"errors"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	ERR_SERVICE_NOT_UNDER_REVIEW = "service not under review"
)

type Service struct {
	serviceStore    service_repo.ServiceRepository
	vendorStore     vendor_repo.VendorRepository
	notifStore      notification_repo.NotificationRepository
	resourceService *resource_service.Service
	ws              websocket.Socket
	encrypt         core.Encryption
}

func NewService(
	serviceStore service_repo.ServiceRepository,
	vendorStore vendor_repo.VendorRepository,
	notifStore notification_repo.NotificationRepository,
	resourceService *resource_service.Service,
	ws websocket.Socket,
	encrypt core.Encryption,
) *Service {
	return &Service{
		serviceStore:    serviceStore,
		vendorStore:     vendorStore,
		notifStore:      notifStore,
		resourceService: resourceService,
		ws:              ws,
		encrypt:         encrypt,
	}
}

func (s *Service) GetPendingServicesList(limit, offset int) ([]dto.PendingServiceListItem, error) {
	services, err := s.serviceStore.GetAllUnderReview(limit, offset)
	if err != nil {
		return nil, err
	}

	pendingServices := make([]dto.PendingServiceListItem, 0)
	for _, service := range services {
		vendor, err := s.vendorStore.FindById(service.VendorId)
		if err != nil {
			return nil, err
		}

		pendingServices = append(pendingServices, dto.PendingServiceListItem{
			ID:          service.Id,
			VendorName:  utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
			VendorEmail: utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
			CreatedAt:   utils.FormatDate(service.CreatedAt),
		})
	}

	return pendingServices, nil
}

func (s *Service) PendingServiceDetail(serviceId string) (*dto.PendingService, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	vendor, err := s.vendorStore.FindById(service.VendorId)
	if err != nil {
		return nil, err
	}

	pendingService := &dto.PendingService{
		Vendor: dto.Vendor{
			Id:       vendor.VendorId,
			Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
			ImageURL: vendor.User.ImageUrl,
			Address:  utils.Try(s.encrypt.DecryptString(vendor.User.Address.Address)),
			Phone:    utils.Try(s.encrypt.DecryptString(vendor.User.Phone)),
			Socials: slices.AppendSeq(
				make([]dto.Social, 0),
				utils.Map(vendor.User.Socials, func(social models.SocialModel) dto.Social {
					return dto.Social{
						Id:    social.Id,
						Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
						Title: utils.Must(s.encrypt.DecryptString(social.Title)),
						URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
					}
				}),
			),
			Expertise: slices.AppendSeq(
				make([]dto.Expertise, 0),
				utils.Map(vendor.Expertise, func(expertise models.ExpertiseModel) dto.Expertise {
					return dto.Expertise{
						Title:              expertise.Title,
						DateApplied:        utils.FormatDate(expertise.DateApplied),
						DateApproved:       utils.FormatDate(expertise.DateApproved.String),
						SupportingDocument: utils.Must(s.resourceService.SignURLWithDefaultDuration(expertise.SupportingImageUrl)),
					}
				}),
			),
			Identification: dto.Identification{
				Type:          vendor.User.Identification.Type,
				IdNumber:      utils.Must(s.encrypt.DecryptString(vendor.User.Identification.ReferenceNumber)),
				FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.FrontImageUrl)),
				BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.BackImageUrl)),
			},
			Rating:       vendor.Rating,
			JoinedAt:     utils.FormatDate(vendor.JoinedAt),
			DateVerified: utils.FormatDate(vendor.User.VerifiedAt.String),
			IsRestricted: vendor.User.Restricted,
			IsBanned:     vendor.User.Banned,
		},
		Service: dto.Service{
			Id:          service.Id,
			VendorId:    service.VendorId,
			Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
			Description: utils.Must(s.encrypt.DecryptString(service.Description)),
			Price:       service.Price,
			PricingType: string(service.PricingType),
			Tags: slices.AppendSeq(
				make([]string, 0),
				utils.Map(service.Tags, func(t *models.TagModel) string { return t.Title }),
			),
			Extras: slices.AppendSeq(
				make([]dto.Extra, 0),
				utils.Map(service.Extras, func(x *models.ExtraModel) dto.Extra {
					return dto.Extra{
						Id:          x.Id,
						Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
						Description: utils.Must(s.encrypt.DecryptString(x.Description)),
						Price:       x.Price,
					}
				}),
			),
			Images: slices.AppendSeq(
				make([]dto.Image, 0),
				utils.Map(service.Images, func(img *models.ServicePhotoModel) dto.Image {
					return dto.Image{Id: img.Id, URL: img.Url}
				}),
			),
			Address: dto.Address{
				Address:   utils.Must(s.encrypt.DecryptString(service.Address.Address)),
				Latitude:  service.Address.Latitude,
				Longitude: service.Address.Longitude,
			},
			CreatedAt:    utils.FormatDate(service.CreatedAt),
			UpdatedAt:    utils.FormatDate(service.UpdatedAt),
			Status:       string(service.Status),
			RejectReason: utils.Must(s.encrypt.DecryptString(service.RejectReason.String)),
			AcceptedAt:   utils.FormatDate(service.AcceptedAt.String),
			RejectedAt:   utils.FormatDate(service.RejectedAt.String),
		},
	}

	return pendingService, nil
}

func (s *Service) Accept(serviceId string) error {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return err
	}

	if service.Status != models.SERVICE_STATUS_UNDER_REVIEW {
		return errors.New(ERR_SERVICE_NOT_UNDER_REVIEW)
	}

	if err := s.serviceStore.Accept(service.Id); err != nil {
		return err
	}

	vendor, err := s.vendorStore.FindById(service.VendorId)
	if err != nil {
		return err
	}

	notificationHeading := "Service accepted"
	notificationContent := "Congratulations! Your service passed the review."

	notification := &models.NotificationModel{
		Recipient: vendor.VendorId,
		Type:      "success",
		Title:     "Service Submission Accepted",
		Content: fmt.Sprintf(
			"Congratulations! The admin accepted your submitted service titled: %s",
			utils.Must(s.encrypt.DecryptString(service.Title)),
		),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: vendor.VendorId,
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
	if err := oneSignal.NewUrgentNotification(vendor.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: vendor.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	serviceEvt := &websocket.EventModel{
		ReceiverId: vendor.VendorId,
		Type:       websocket.EVT_SERVICE_ACCEPTED,
		Payload:    utils.Mapper{"id": service.Id},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(serviceEvt)

	return nil
}

func (s *Service) Reject(serviceId, reason string) error {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return err
	}

	if service.Status != models.SERVICE_STATUS_UNDER_REVIEW {
		return errors.New(ERR_SERVICE_NOT_UNDER_REVIEW)
	}

	encryptedReason := utils.Must(s.encrypt.EncryptString(reason))
	if err := s.serviceStore.Reject(service.Id, encryptedReason); err != nil {
		return err
	}

	vendor, err := s.vendorStore.FindById(service.VendorId)
	if err != nil {
		return err
	}

	notificationHeading := "Service rejected"
	notificationContent := "Admin rejected service submission"

	notification := &models.NotificationModel{
		Recipient: vendor.VendorId,
		Type:      "fail",
		Title:     "Service Submission Rejected",
		Content: fmt.Sprintf(
			"Unfortunately, your submitted service title: \"%s\" did not meet platform guidelines and was rejected. Reason for rejection: %s",
			utils.Must(s.encrypt.DecryptString(service.Title)),
			reason,
		),
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: vendor.VendorId,
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
	if err := oneSignal.NewUrgentNotification(vendor.VendorId, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	notifEvent := &websocket.EventModel{
		ReceiverId: vendor.VendorId,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	serviceEvt := &websocket.EventModel{
		ReceiverId: vendor.VendorId,
		Type:       websocket.EVT_SERVICE_REJECTED,
		Payload: utils.Mapper{
			"id":     service.Id,
			"reason": reason,
		},
	}

	s.ws.Send(notifEvent)
	s.ws.Send(serviceEvt)

	return nil
}
