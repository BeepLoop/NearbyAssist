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
			Identification: dto.Identification{
				Type:          identification.Type,
				IdNumber:      utils.Must(s.encrypt.DecryptString(identification.ReferenceNumber)),
				FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.FrontImageUrl)),
				BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(identification.BackImageUrl)),
			},
			CreatedAt:    utils.FormatDate(user.CreatedAt),
			DateVerified: utils.FormatDate(user.VerifiedAt.String),
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
						Rate:        utils.StringToFloat64ElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloat64ElseZero(booking.Cost),
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
						Rate:        utils.StringToFloat64ElseZero(service.Rate),
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Cost:   utils.StringToFloat64ElseZero(h.Cost),
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
