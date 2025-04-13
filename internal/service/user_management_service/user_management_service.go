package user_management_service

import (
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	notification_repo "nearbyassist/internal/repository/notification"
	service_repo "nearbyassist/internal/repository/service"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
	"time"
)

type Service struct {
	userStore       user_repo.UserRepository
	vendorStore     vendor_repo.VendorRepository
	notifStore      notification_repo.NotificationRepository
	bookingStore    booking_repo.BookingRepository
	serviceStore    service_repo.ServiceRepository
	resourceService *resource_service.Service
	ws              websocket.Socket
	encrypt         core.Encryption
	hash            core.Hash
}

func NewService(
	userStore user_repo.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	notifStore notification_repo.NotificationRepository,
	bookingStore booking_repo.BookingRepository,
	serviceStore service_repo.ServiceRepository,
	resourceService *resource_service.Service,
	ws websocket.Socket,
	encrypt core.Encryption,
	hash core.Hash,
) *Service {
	return &Service{
		userStore:       userStore,
		vendorStore:     vendorStore,
		notifStore:      notifStore,
		bookingStore:    bookingStore,
		serviceStore:    serviceStore,
		resourceService: resourceService,
		ws:              ws,
		encrypt:         encrypt,
		hash:            hash,
	}
}

func (s *Service) GetUserAccountDetail(userId string) (*dto.UserAccountDetail, error) {
	user, err := s.userStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	identification, err := s.userStore.GetIdentification(user.Id)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetConfirmed(userId, "client")
	if err != nil {
		return nil, err
	}

	history, err := s.bookingStore.GetHistory(userId, "client")
	if err != nil {
		return nil, err
	}

	data := &dto.UserAccountDetail{
		User: dto.User{
			Id:       user.Id,
			Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
			ImageURL: user.ImageUrl,
			Address:  utils.Try(s.encrypt.DecryptString(user.Address.String)),
			Phone:    utils.Try(s.encrypt.DecryptString(user.Phone.String)),
			Socials: slices.AppendSeq(
				make([]string, 0),
				utils.Map(user.Socials, func(social string) string {
					return utils.Must(s.encrypt.DecryptString(social))
				}),
			),
			Identification: dto.Identification{
				Type:          identification.Type,
				IdNumber:      utils.Must(s.encrypt.DecryptString(identification.IdNumber)),
				FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.FrontImage)),
				BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.BackImage)),
			},
			CreatedAt:    utils.FormatDate(user.CreatedAt),
			DateVerified: utils.FormatDate(user.VerifiedAt),
			IsRestricted: user.Restricted,
			IsBanned:     user.Banned,
		},
		ActiveBookings: slices.AppendSeq(
			make([]dto.Booking, 0),
			utils.Map(bookings, func(booking *models.BookingModel) dto.Booking {
				vendor, _ := s.vendorStore.FindById(booking.VendorId)
				service, _ := s.serviceStore.FindById(booking.ServiceId)

				return dto.Booking{
					Id: booking.Id,
					Client: dto.User{
						Id:       user.Id,
						Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
						ImageURL: user.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(vendor.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(vendor.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(vendor.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						CreatedAt:    utils.FormatDate(user.CreatedAt),
						DateVerified: utils.FormatDate(user.VerifiedAt),
						IsRestricted: user.Restricted,
						IsBanned:     user.Banned,
					},
					Vendor: dto.Vendor{
						Id:       vendor.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(vendor.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(vendor.Email)),
						ImageURL: vendor.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(vendor.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(vendor.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(vendor.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						Rating:       vendor.Rating,
						JoinedAt:     utils.FormatDate(vendor.JoinedAt),
						DateVerified: utils.FormatDate(vendor.VerifiedAt),
					},
					Service: dto.Service{
						Id:          service.Id,
						Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
						Description: utils.Must(s.encrypt.DecryptString(service.Description)),
						Rate:        utils.StringToFloatElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloatElseZero(booking.Cost),
					Status: string(booking.Status),
					Extras: slices.AppendSeq(
						make([]dto.Extra, 0),
						utils.Map(booking.Extras, func(x *models.ExtraModel) dto.Extra {
							return dto.Extra{
								Id:          x.Id,
								Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
								Description: utils.Must(s.encrypt.DecryptString(x.Description)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:    utils.FormatDate(booking.CreatedAt),
					ScheduledAt:  utils.FormatDate(booking.ScheduledAt.String),
					UpdatedAt:    utils.FormatDate(booking.UpdatedAt),
					CancelReason: utils.Try(s.encrypt.DecryptString(booking.CancelReason.String)),
				}
			}),
		),
		History: slices.AppendSeq(
			make([]dto.Booking, 0),
			utils.Map(history, func(h *models.BookingModel) dto.Booking {
				vendor, _ := s.vendorStore.FindById(h.VendorId)
				service, _ := s.serviceStore.FindById(h.ServiceId)

				return dto.Booking{
					Id: h.Id,
					Client: dto.User{
						Id:       user.Id,
						Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
						ImageURL: user.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(user.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(user.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(user.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						CreatedAt:    utils.FormatDate(user.CreatedAt),
						DateVerified: utils.FormatDate(user.VerifiedAt),
						IsRestricted: user.Restricted,
						IsBanned:     user.Banned,
					},
					Vendor: dto.Vendor{
						Id:       vendor.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(vendor.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(vendor.Email)),
						ImageURL: vendor.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(vendor.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(vendor.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(vendor.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						Rating:       vendor.Rating,
						JoinedAt:     utils.FormatDate(vendor.JoinedAt),
						DateVerified: utils.FormatDate(vendor.VerifiedAt),
					},
					Service: dto.Service{
						Id:          service.Id,
						Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
						Description: utils.Must(s.encrypt.DecryptString(service.Description)),
						Rate:        utils.StringToFloatElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloatElseZero(h.Cost),
					Status: string(h.Status),
					Extras: slices.AppendSeq(
						make([]dto.Extra, 0),
						utils.Map(h.Extras, func(x *models.ExtraModel) dto.Extra {
							return dto.Extra{
								Id:          x.Id,
								Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
								Description: utils.Must(s.encrypt.DecryptString(x.Description)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:    utils.FormatDate(h.CreatedAt),
					ScheduledAt:  utils.FormatDate(h.ScheduledAt.String),
					UpdatedAt:    utils.FormatDate(h.UpdatedAt),
					CancelReason: utils.Try(s.encrypt.DecryptString(h.CancelReason.String)),
				}
			}),
		),
	}

	return data, nil
}

func (s *Service) GetVendorAccountDetail(userId string) (*dto.VendorAccountDetail, error) {
	account, err := s.vendorStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	identification, err := s.userStore.GetIdentification(account.VendorId)
	if err != nil {
		return nil, err
	}

	expertises, err := s.vendorStore.GetAllExpertise(account.VendorId)
	if err != nil {
		return nil, err
	}

	services, err := s.vendorStore.GetVendorServiceList(userId)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingStore.GetConfirmed(userId, "vendor")
	if err != nil {
		return nil, err
	}

	history, err := s.bookingStore.GetHistory(userId, "vendor")
	if err != nil {
		return nil, err
	}

	data := &dto.VendorAccountDetail{
		Vendor: dto.Vendor{
			Id:       account.VendorId,
			Name:     utils.Must(s.encrypt.DecryptString(account.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(account.Email)),
			ImageURL: account.ImageUrl,
			Address:  utils.Try(s.encrypt.DecryptString(account.Address.String)),
			Phone:    utils.Try(s.encrypt.DecryptString(account.Phone.String)),
			Socials: slices.AppendSeq(
				make([]string, 0),
				utils.Map(account.Socials, func(social string) string {
					return utils.Must(s.encrypt.DecryptString(social))
				}),
			),
			Identification: dto.Identification{
				Type:          identification.Type,
				IdNumber:      utils.Must(s.encrypt.DecryptString(identification.IdNumber)),
				FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.FrontImage)),
				BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.BackImage)),
			},
			Rating: account.Rating,
			Expertise: slices.AppendSeq(
				make([]dto.Expertise, 0),
				utils.Map(expertises, func(e *models.UserExpertiseModel) dto.Expertise {
					return dto.Expertise{
						Title:              e.Expertise,
						DateApplied:        utils.FormatDate(e.DateApplied),
						DateApproved:       utils.FormatDate(e.DateApproved),
						SupportingDocument: utils.Must(s.resourceService.SignURLWithDefaultDuration(e.SupportingDocumentImage)),
					}
				}),
			),
			JoinedAt:     utils.FormatDate(account.JoinedAt),
			DateVerified: utils.FormatDate(account.VerifiedAt),
			IsRestricted: account.Restricted,
			IsBanned:     account.Banned,
		},
		Services: slices.AppendSeq(
			make([]dto.Service, 0),
			utils.Map(services, func(service *models.ServiceModel) dto.Service {
				images, _ := s.serviceStore.GetPhotos(service.Id)

				return dto.Service{
					Id:          service.Id,
					Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
					Description: utils.Must(s.encrypt.DecryptString(service.Description)),
					Rate:        utils.StringToFloatElseZero(service.Rate),
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
						utils.Map(images, func(img *models.ServicePhotoModel) dto.Image {
							return dto.Image{Id: img.Id, URL: img.Url}
						}),
					),
					CreatedAt: utils.FormatDate(service.CreatedAt),
					UpdatedAt: utils.FormatDate(service.UpdatedAt),
				}
			}),
		),
		ActiveBookings: slices.AppendSeq(
			make([]dto.Booking, 0),
			utils.Map(bookings, func(booking *models.BookingModel) dto.Booking {
				client, _ := s.userStore.FindById(booking.ClientId)
				service, _ := s.serviceStore.FindById(booking.ServiceId)

				return dto.Booking{
					Id: booking.Id,
					Client: dto.User{
						Id:       client.Id,
						Name:     utils.Must(s.encrypt.DecryptString(client.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(client.Email)),
						ImageURL: client.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(client.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(client.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(client.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						CreatedAt:    utils.FormatDate(client.CreatedAt),
						DateVerified: utils.FormatDate(client.VerifiedAt),
						IsRestricted: client.Restricted,
						IsBanned:     client.Banned,
					},
					Vendor: dto.Vendor{
						Id:       account.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(account.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(account.Email)),
						ImageURL: account.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(account.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(account.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(account.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						Rating:       account.Rating,
						JoinedAt:     utils.FormatDate(account.JoinedAt),
						DateVerified: utils.FormatDate(account.VerifiedAt),
						IsRestricted: account.Restricted,
						IsBanned:     account.Banned,
					},
					Service: dto.Service{
						Id:          service.Id,
						Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
						Description: utils.Must(s.encrypt.DecryptString(service.Description)),
						Rate:        utils.StringToFloatElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloatElseZero(booking.Cost),
					Status: string(booking.Status),
					Extras: slices.AppendSeq(
						make([]dto.Extra, 0),
						utils.Map(booking.Extras, func(x *models.ExtraModel) dto.Extra {
							return dto.Extra{
								Id:          x.Id,
								Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
								Description: utils.Must(s.encrypt.DecryptString(x.Description)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:    utils.FormatDate(booking.CreatedAt),
					ScheduledAt:  utils.FormatDate(booking.ScheduledAt.String),
					UpdatedAt:    utils.FormatDate(booking.UpdatedAt),
					CancelReason: utils.Try(s.encrypt.DecryptString(booking.CancelReason.String)),
				}
			}),
		),
		History: slices.AppendSeq(
			make([]dto.Booking, 0),
			utils.Map(history, func(h *models.BookingModel) dto.Booking {
				client, _ := s.userStore.FindById(h.ClientId)
				service, _ := s.serviceStore.FindById(h.ServiceId)

				return dto.Booking{
					Id: h.Id,
					Client: dto.User{
						Id:       client.Id,
						Name:     utils.Must(s.encrypt.DecryptString(client.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(client.Email)),
						ImageURL: client.ImageUrl,
						Address:  client.Address.String,
						Phone:    client.Phone.String,
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(client.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						CreatedAt:    utils.FormatDate(client.CreatedAt),
						DateVerified: utils.FormatDate(client.VerifiedAt),
						IsRestricted: client.Restricted,
						IsBanned:     client.Banned,
					},
					Vendor: dto.Vendor{
						Id:       account.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(account.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(account.Email)),
						ImageURL: account.ImageUrl,
						Address:  utils.Try(s.encrypt.DecryptString(account.Address.String)),
						Phone:    utils.Try(s.encrypt.DecryptString(account.Phone.String)),
						Socials: slices.AppendSeq(
							make([]string, 0),
							utils.Map(account.Socials, func(social string) string {
								return utils.Must(s.encrypt.DecryptString(social))
							}),
						),
						Rating:       account.Rating,
						JoinedAt:     utils.FormatDate(account.JoinedAt),
						DateVerified: utils.FormatDate(account.VerifiedAt),
						IsRestricted: account.Restricted,
						IsBanned:     account.Banned,
					},
					Service: dto.Service{
						Id:          service.Id,
						Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
						Description: utils.Must(s.encrypt.DecryptString(service.Description)),
						Rate:        utils.StringToFloatElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloatElseZero(h.Cost),
					Status: string(h.Status),
					Extras: slices.AppendSeq(
						make([]dto.Extra, 0),
						utils.Map(h.Extras, func(x *models.ExtraModel) dto.Extra {
							return dto.Extra{
								Id:          x.Id,
								Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
								Description: utils.Must(s.encrypt.DecryptString(x.Description)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:    utils.FormatDate(h.CreatedAt),
					ScheduledAt:  utils.FormatDate(h.ScheduledAt.String),
					UpdatedAt:    utils.FormatDate(h.UpdatedAt),
					CancelReason: utils.Try(s.encrypt.DecryptString(h.CancelReason.String)),
				}
			}),
		),
	}

	return data, nil
}

func (s *Service) BanUser(userId string) error {
	if err := s.userStore.BanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

	return nil
}

func (s *Service) UnbanUser(userId string) error {
	if err := s.userStore.UnbanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

	return nil
}

func (s *Service) RestrictUser(userId, reason, duration string) error {
	d, err := utils.ParseStringDuration(duration)
	if err != nil {
		return err
	}

	endDate := time.Now().Add(d)
	data := &models.RestrictionModel{
		UserId:  userId,
		Reason:  reason,
		EndTime: utils.FormatDateTime(endDate),
	}

	if err := s.userStore.RestrictUser(data); err != nil {
		return err
	}

	stringifiedDuration := utils.FormatDurationToString(d)

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "generic",
		Title:     "Account Restricted " + stringifiedDuration,
		Content:   reason,
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

	notificationHeading := "Account Restricted!"
	notificationContent := "You commited a violation resulting to account restriction."

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

func (s *Service) UnrestrictUser(userId string) error {
	if err := s.userStore.ForceLiftRestriction(userId); err != nil {
		return err
	}

	notification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     "Restriction Lifted",
		Content:   "The restriction to your account has been lifted by the administrator. Avoid committing violations to prevent future restrictions.",
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: userId,
		Type:      "success",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	notificationHeading := "Account Restriction Lifted!"
	notificationContent := "Your account restriction has been lifted."

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
