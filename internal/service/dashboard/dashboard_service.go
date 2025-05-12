package dashboard_service

import (
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	service_repo "nearbyassist/internal/repository/service"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	dashboardStore dashboard_repo.DashboardRepository
	serviceStore   service_repo.ServiceRepository
	bookingStore   booking_repo.BookingRepository
	userStore      user_repo.UserRepository
	vendorStore    vendor_repo.VendorRepository
	encrypt        core.Encryption
}

func NewService(
	dashboardStore dashboard_repo.DashboardRepository,
	serviceStore service_repo.ServiceRepository,
	bookingStore booking_repo.BookingRepository,
	userStore user_repo.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	encrypt core.Encryption,
) *Service {
	return &Service{
		dashboardStore: dashboardStore,
		serviceStore:   serviceStore,
		bookingStore:   bookingStore,
		userStore:      userStore,
		vendorStore:    vendorStore,
		encrypt:        encrypt,
	}
}

func (s *Service) GetDashbaordData() (*dto.Dashboard, error) {
	recentUsers, err := s.dashboardStore.RecentUsers()
	if err != nil {
		return nil, err
	}

	weeklyBookingStats := utils.Must(s.dashboardStore.GetBookingData())
	bookings, err := s.dashboardStore.GetBookingsThisWeek()
	if err != nil {
		fmt.Println("error retrieving booking this week: ", err.Error())
		return nil, err
	}

	weeklyBookingData := dto.WeeklyBookingData{
		Total:      weeklyBookingStats.Total,
		Daily:      weeklyBookingStats.Daily,
		Difference: weeklyBookingStats.Difference,
		Bookings: slices.AppendSeq(
			make([]dto.Booking, 0),
			utils.Map(bookings, func(booking *models.BookingModel) dto.Booking {
				client, err := s.userStore.FindById(booking.ClientId)
				if err != nil {
					fmt.Println("error retrieving client in booking: ", err.Error())
				}

				vendor, err := s.vendorStore.FindById(booking.VendorId)
				if err != nil {
					fmt.Println("error retrieving vendor in booking: ", err.Error())
				}

				service, err := s.serviceStore.FindById(booking.ServiceId)
				if err != nil {
					fmt.Println("error retrieving service in booking: ", err.Error())
				}

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
						CreatedAt: utils.FormatDate(service.CreatedAt),
						UpdatedAt: utils.FormatDate(service.UpdatedAt),
					},
					Quantity: booking.Quantity,
					Cost:     booking.Cost,
					Status:   string(booking.Status),
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
					CreatedAt:     utils.FormatDate(booking.CreatedAt),
					ScheduleStart: utils.FormatDate(booking.ScheduleStart.String),
					ScheduleEnd:   utils.FormatDate(booking.ScheduleEnd.String),
					UpdatedAt:     utils.FormatDate(booking.UpdatedAt),
					CancelReason:  utils.Try(s.encrypt.DecryptString(booking.CancelReason.String)),
				}
			}),
		),
	}

	data := &dto.Dashboard{
		Users: dto.DashboardUsers{
			Total:      utils.Must(s.dashboardStore.TotalUsers()),
			Reported:   utils.Must(s.dashboardStore.TotalReported()),
			Restricted: utils.Must(s.dashboardStore.TotalRestricted()),
			Vendor:     utils.Must(s.dashboardStore.TotalVendors()),
			Recent: slices.AppendSeq(
				make([]dto.User, 0),
				utils.Map(recentUsers, func(user *models.UserModel) dto.User {
					return dto.User{
						Id:        user.Id,
						Name:      utils.Must(s.encrypt.DecryptString(user.Name)),
						Email:     utils.Must(s.encrypt.DecryptString(user.Email)),
						ImageURL:  user.ImageUrl,
						CreatedAt: utils.DateMonth(user.CreatedAt),
					}
				}),
			),
		},
		Services: dto.DashboardServices{
			Total:  utils.Must(s.dashboardStore.TotalServices()),
			Active: utils.Must(s.dashboardStore.TotalActiveServices()),
		},
		WeeklyBooking: weeklyBookingData,
		SearchTrend: dto.Trend{
			Searches: searchhistory.New().GetAll(),
		},
		Application: dto.Application{
			Pending: utils.Must(s.dashboardStore.TotalPendingApplications()),
		},
		Report: dto.DashboardReports{
			Pending: utils.Must(s.dashboardStore.TotalActiveReports()),
		},
	}

	return data, nil
}
