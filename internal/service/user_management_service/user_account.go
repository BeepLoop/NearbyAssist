package user_management_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
)

func (s *Service) GetUserAccountDetail(userId string) (*dto.UserAccountDetail, error) {
	user, err := s.userStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	identification := dto.Identification{}
	if user.HasSubmittedIdentification {
		identification = dto.Identification{
			Type:          user.Identification.Type,
			IdNumber:      utils.Must(s.encrypt.DecryptString(user.Identification.ReferenceNumber)),
			FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.FrontImageUrl)),
			BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.BackImageUrl)),
		}
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
			Address:  utils.Try(s.encrypt.DecryptString(user.Address.Address)),
			Phone:    utils.Try(s.encrypt.DecryptString(user.Phone)),
			Socials: slices.AppendSeq(
				make([]dto.Social, 0),
				utils.Map(user.Socials, func(social models.SocialModel) dto.Social {
					return dto.Social{
						Id:    social.Id,
						Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
						Title: utils.Must(s.encrypt.DecryptString(social.Title)),
						URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
					}
				}),
			),
			Identification:             identification,
			CreatedAt:                  utils.FormatDate(user.CreatedAt),
			DateVerified:               utils.FormatDate(user.VerifiedAt.String),
			IsRestricted:               user.Restricted,
			IsBanned:                   user.Banned,
			HasSubmittedIdentification: user.HasSubmittedIdentification,
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
						CreatedAt:    utils.FormatDate(user.CreatedAt),
						DateVerified: utils.FormatDate(user.VerifiedAt.String),
						IsRestricted: user.Restricted,
						IsBanned:     user.Banned,
					},
					Vendor: dto.Vendor{
						Id:       vendor.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
						ImageURL: vendor.User.ImageUrl,
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
						Rating:       vendor.Rating,
						JoinedAt:     utils.FormatDate(vendor.JoinedAt),
						DateVerified: utils.FormatDate(vendor.User.VerifiedAt.String),
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
					ServiceId:          booking.ServiceId,
					ServiceTitle:       utils.Must(s.encrypt.DecryptString(booking.ServiceTitle)),
					ServiceDescription: utils.Must(s.encrypt.DecryptString(booking.ServiceDescription)),
					Price:              booking.Price,
					PricingType:        string(booking.PricingType),
					Quantity:           booking.Quantity,
					Cost:               booking.Cost,
					Status:             string(booking.Status),
					Extras: slices.AppendSeq(
						make([]dto.BookingExtra, 0),
						utils.Map(booking.Extras, func(x *models.BookingExtraModel) dto.BookingExtra {
							return dto.BookingExtra{
								BookingId:   x.BookingId,
								Title:       utils.Must(s.encrypt.DecryptString(x.ExtraTitle)),
								Description: utils.Must(s.encrypt.DecryptString(x.ExtraDescription)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:     utils.FormatDate(booking.CreatedAt),
					ScheduleStart: utils.FormatDate(booking.ScheduleStart.String),
					ScheduleEnd:   utils.FormatDate(booking.ScheduleEnd.String),
					UpdatedAt:     utils.FormatDate(booking.UpdatedAt),
					CancelReason:  utils.Try(s.encrypt.DecryptString(booking.CancelReason.String)),
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
						Socials: slices.AppendSeq(
							make([]dto.Social, 0),
							utils.Map(user.Socials, func(social models.SocialModel) dto.Social {
								return dto.Social{
									Id:    social.Id,
									Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
									Title: utils.Must(s.encrypt.DecryptString(social.Title)),
									URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
								}
							}),
						),
						CreatedAt:    utils.FormatDate(user.CreatedAt),
						DateVerified: utils.FormatDate(user.VerifiedAt.String),
						IsRestricted: user.Restricted,
						IsBanned:     user.Banned,
					},
					Vendor: dto.Vendor{
						Id:       vendor.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
						ImageURL: vendor.User.ImageUrl,
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
						Rating:       vendor.Rating,
						JoinedAt:     utils.FormatDate(vendor.JoinedAt),
						DateVerified: utils.FormatDate(vendor.User.VerifiedAt.String),
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
					ServiceId:          h.ServiceId,
					ServiceTitle:       utils.Must(s.encrypt.DecryptString(h.ServiceTitle)),
					ServiceDescription: utils.Must(s.encrypt.DecryptString(h.ServiceDescription)),
					Price:              h.Price,
					PricingType:        string(h.PricingType),
					Quantity:           h.Quantity,
					Cost:               h.Cost,
					Status:             string(h.Status),
					Extras: slices.AppendSeq(
						make([]dto.BookingExtra, 0),
						utils.Map(h.Extras, func(x *models.BookingExtraModel) dto.BookingExtra {
							return dto.BookingExtra{
								BookingId:   x.BookingId,
								Title:       utils.Must(s.encrypt.DecryptString(x.ExtraTitle)),
								Description: utils.Must(s.encrypt.DecryptString(x.ExtraDescription)),
								Price:       x.Price,
							}
						}),
					),
					CreatedAt:     utils.FormatDate(h.CreatedAt),
					ScheduleStart: utils.FormatDate(h.ScheduleStart.String),
					ScheduleEnd:   utils.FormatDate(h.ScheduleEnd.String),
					UpdatedAt:     utils.FormatDate(h.UpdatedAt),
					CancelReason:  utils.Try(s.encrypt.DecryptString(h.CancelReason.String)),
				}
			}),
		),
	}

	return data, nil
}
