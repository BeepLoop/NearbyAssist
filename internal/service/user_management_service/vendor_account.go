package user_management_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
)

func (s *Service) GetVendorAccountDetail(userId string) (*dto.VendorAccountDetail, error) {
	account, err := s.vendorStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	identification, err := s.userStore.GetIdentification(account.VendorId)
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
			Name:     utils.Must(s.encrypt.DecryptString(account.User.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(account.User.Email)),
			ImageURL: account.User.ImageUrl,
			Address:  utils.Try(s.encrypt.DecryptString(account.User.Address.Address)),
			Phone:    utils.Try(s.encrypt.DecryptString(account.User.Phone)),
			Socials: slices.AppendSeq(
				make([]dto.Social, 0),
				utils.Map(account.User.Socials, func(social models.SocialModel) dto.Social {
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
			Rating: account.Rating,
			Expertise: slices.AppendSeq(
				make([]dto.Expertise, 0),
				utils.Map(account.Expertise, func(e models.ExpertiseModel) dto.Expertise {
					return dto.Expertise{
						Title:              e.Title,
						DateApplied:        utils.FormatDate(e.DateApplied),
						DateApproved:       utils.FormatDate(e.DateApproved.String),
						SupportingDocument: utils.Must(s.resourceService.SignURLWithDefaultDuration(e.SupportingImageUrl)),
					}
				}),
			),
			JoinedAt:     utils.FormatDate(account.JoinedAt),
			DateVerified: utils.FormatDate(account.User.VerifiedAt.String),
			IsRestricted: account.User.Restricted,
			IsBanned:     account.User.Banned,
		},
		Services: slices.AppendSeq(
			make([]dto.Service, 0),
			utils.Map(services, func(service *models.ServiceModel) dto.Service {
				return dto.Service{
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
						Socials: slices.AppendSeq(
							make([]dto.Social, 0),
							utils.Map(client.Socials, func(social models.SocialModel) dto.Social {
								return dto.Social{
									Id:    social.Id,
									Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
									Title: utils.Must(s.encrypt.DecryptString(social.Title)),
									URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
								}
							}),
						),
						CreatedAt:    utils.FormatDate(client.CreatedAt),
						DateVerified: utils.FormatDate(client.VerifiedAt.String),
						IsRestricted: client.Restricted,
						IsBanned:     client.Banned,
					},
					Vendor: dto.Vendor{
						Id:       account.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(account.User.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(account.User.Email)),
						ImageURL: account.User.ImageUrl,
						Socials: slices.AppendSeq(
							make([]dto.Social, 0),
							utils.Map(account.User.Socials, func(social models.SocialModel) dto.Social {
								return dto.Social{
									Id:    social.Id,
									Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
									Title: utils.Must(s.encrypt.DecryptString(social.Title)),
									URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
								}
							}),
						),
						Rating:       account.Rating,
						JoinedAt:     utils.FormatDate(account.JoinedAt),
						DateVerified: utils.FormatDate(account.User.VerifiedAt.String),
						IsRestricted: account.User.Restricted,
						IsBanned:     account.User.Banned,
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
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
				client, _ := s.userStore.FindById(h.ClientId)
				service, _ := s.serviceStore.FindById(h.ServiceId)

				return dto.Booking{
					Id: h.Id,
					Client: dto.User{
						Id:       client.Id,
						Name:     utils.Must(s.encrypt.DecryptString(client.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(client.Email)),
						ImageURL: client.ImageUrl,
						Socials: slices.AppendSeq(
							make([]dto.Social, 0),
							utils.Map(client.Socials, func(social models.SocialModel) dto.Social {
								return dto.Social{
									Id:    social.Id,
									Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
									Title: utils.Must(s.encrypt.DecryptString(social.Title)),
									URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
								}
							}),
						),
						CreatedAt:    utils.FormatDate(client.CreatedAt),
						DateVerified: utils.FormatDate(client.VerifiedAt.String),
						IsRestricted: client.Restricted,
						IsBanned:     client.Banned,
					},
					Vendor: dto.Vendor{
						Id:       account.VendorId,
						Name:     utils.Must(s.encrypt.DecryptString(account.User.Name)),
						Email:    utils.Must(s.encrypt.DecryptString(account.User.Email)),
						ImageURL: account.User.ImageUrl,
						Socials: slices.AppendSeq(
							make([]dto.Social, 0),
							utils.Map(account.User.Socials, func(social models.SocialModel) dto.Social {
								return dto.Social{
									Id:    social.Id,
									Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
									Title: utils.Must(s.encrypt.DecryptString(social.Title)),
									URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
								}
							}),
						),
						Rating:       account.Rating,
						JoinedAt:     utils.FormatDate(account.JoinedAt),
						DateVerified: utils.FormatDate(account.User.VerifiedAt.String),
						IsRestricted: account.User.Restricted,
						IsBanned:     account.User.Banned,
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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
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
