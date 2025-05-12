package complaint_service

import (
	"errors"
	"mime/multipart"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	admin_repo "nearbyassist/internal/repository/admin"
	booking_repo "nearbyassist/internal/repository/booking"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	notification_repo "nearbyassist/internal/repository/notification"
	report_user_repo "nearbyassist/internal/repository/report_user"
	service_repo "nearbyassist/internal/repository/service"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"slices"
	"strconv"
)

type Service struct {
	reportUserStore report_user_repo.ReportUserRepository
	adminStore      admin_repo.AdminRepository
	userStore       user_repo.UserRepository
	vendorStore     vendor_repo.VendorRepository
	bookingStore    booking_repo.BookingRepository
	serviceStore    service_repo.ServiceRepository
	bugReportStore  bug_report_repo.BugReportRepository
	notifStore      notification_repo.NotificationRepository
	resourceService *resource_service.Service
	ws              websocket.Socket
	fs              fs.FileStorage
	encrypt         core.Encryption
	jwt             core.Authenticator
}

func NewService(
	reportUserStore report_user_repo.ReportUserRepository,
	adminStore admin_repo.AdminRepository,
	userStore user_repo.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	bookingStore booking_repo.BookingRepository,
	serviceStore service_repo.ServiceRepository,
	bugReportStore bug_report_repo.BugReportRepository,
	notifStore notification_repo.NotificationRepository,
	resourceService *resource_service.Service,
	ws websocket.Socket,
	fs fs.FileStorage,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		reportUserStore: reportUserStore,
		adminStore:      adminStore,
		userStore:       userStore,
		vendorStore:     vendorStore,
		bookingStore:    bookingStore,
		serviceStore:    serviceStore,
		bugReportStore:  bugReportStore,
		notifStore:      notifStore,
		resourceService: resourceService,
		ws:              ws,
		fs:              fs,
		encrypt:         encrypt,
		jwt:             jwt,
	}
}

func (s *Service) CreateBugReport(req *request.BugReportPayload, files []*multipart.FileHeader) error {
	reportData := &models.BugReportModel{
		Title:  req.Title,
		Detail: req.Detail,
		Images: make([]string, 0),
	}

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return err
		}

		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return err
		}

		fileData := fs.File{
			Data:     cipher,
			Category: fs.BUG_REPORT_DIR,
		}
		url, err := s.fs.SaveFile(fileData)
		if err != nil {
			return err
		}

		reportData.Images = append(reportData.Images, url)
	}

	reportData.Title = utils.Must(s.encrypt.EncryptString(req.Title))
	reportData.Detail = utils.Must(s.encrypt.EncryptString(req.Detail))

	if err := s.bugReportStore.Create(reportData); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetBugReports(limit, offset int) ([]*models.BugReportModel, error) {
	complaints, err := s.bugReportStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, complaint := range complaints {
		complaint.Title = utils.Must(s.encrypt.DecryptString(complaint.Title))
		complaint.Detail = utils.Must(s.encrypt.DecryptString(complaint.Detail))
	}

	return complaints, nil
}

func (s *Service) CompleteBug(bugId string) error {
	id, err := strconv.Atoi(bugId)
	if err != nil {
		return err
	}

	return s.bugReportStore.CompleteBug(id)
}

func (s *Service) ReportUser(bearerToken string, req *request.ReportUserPayload, files []*multipart.FileHeader) error {
	reporterId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	reportData := &models.UserReportModel{
		ReporterUserId: reporterId,
		ReportedUserId: req.UserId,
		Category:       models.UserReportCategory(req.Category),
		BookingIdInput: req.BookingId,
		Reason:         utils.Must(s.encrypt.EncryptString(req.Reason)),
		Detail:         utils.Must(s.encrypt.EncryptString(req.Detail)),
		Images:         make([]string, 0),
	}

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
		}

		fileData := fs.File{
			Data:     utils.Must(s.encrypt.EncryptFile(bytes)),
			Category: fs.REPORT_USER_DIR,
		}
		url, err := s.fs.SaveFile(fileData)
		if err != nil {
			return err
		}

		reportData.Images = append(reportData.Images, url)
	}

	if err := s.reportUserStore.Create(reportData); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetReportList(limit, offset int) ([]dto.ReportItem, error) {
	reports, err := s.reportUserStore.GetAllWithStatus("pending", limit, offset)
	if err != nil {
		return nil, err
	}

	data := slices.AppendSeq(
		make([]dto.ReportItem, 0),
		utils.Map(reports, func(report *models.UserReportModel) dto.ReportItem {
			return dto.ReportItem{
				Id: utils.FormatReportID(report.Id),
				Reported: dto.ReportUser{
					Id:    report.Reported.Id,
					Name:  utils.Must(s.encrypt.DecryptString(report.Reported.Name)),
					Email: utils.Must(s.encrypt.DecryptString(report.Reported.Email)),
				},
				Reporter: dto.ReportUser{
					Id:    report.Reporter.Id,
					Name:  utils.Must(s.encrypt.DecryptString(report.Reporter.Name)),
					Email: utils.Must(s.encrypt.DecryptString(report.Reporter.Email)),
				},
				CreatedAt: utils.FormatDate(report.CreatedAt),
			}
		}),
	)

	return data, nil
}

func (s *Service) GetReport(reportId string) (*dto.UserReportDetail, error) {
	parsedId, err := utils.ParseReportID(reportId)
	if err != nil {
		return nil, err
	}

	report, err := s.reportUserStore.FindById(parsedId)
	if err != nil {
		return nil, err
	}

	reporter, err := s.userStore.FindById(report.ReporterUserId)
	if err != nil {
		return nil, err
	}

	reporterCancelledBooking, err := s.bookingStore.GetHistory(reporter.Id, "client")
	if err != nil {
		return nil, err
	}

	reporterReportsFiled, err := s.reportUserStore.GetAllReportedBy(reporter.Id)
	if err != nil {
		return nil, err
	}

	reported, err := s.userStore.FindById(report.ReportedUserId)
	if err != nil {
		return nil, err
	}

	vendor, err := s.vendorStore.FindById(reported.Id)
	if err != nil {
		return nil, err
	}

	vendorBookings, err := s.vendorStore.GetBookingsWithStatus(vendor.VendorId, "all")
	if err != nil {
		return nil, err
	}

	reportedUserReportHistory, err := s.reportUserStore.GetAllReportedIs(reported.Id)
	if err != nil {
		return nil, err
	}

	reportHistory := slices.AppendSeq(
		make([]dto.PreviousReport, 0),
		utils.Map(
			slices.Collect(utils.Retain(reportedUserReportHistory, func(r *models.UserReportModel) bool {
				return r.Status != models.REPORT_STATUS_PENDING
			})),
			func(report *models.UserReportModel) dto.PreviousReport {
				reporter, _ := s.userStore.FindById(report.ReporterUserId)
				admin, _ := s.adminStore.FindById(report.AdminId.String)

				return dto.PreviousReport{
					Id:               utils.FormatReportID(report.Id),
					ReportedByUserId: report.ReporterUserId,
					ReportedByName:   utils.Must(s.encrypt.DecryptString(reporter.Name)),
					Category:         string(report.Category),
					BookingId:        report.BookingId.String,
					Reason:           utils.Must(s.encrypt.DecryptString(report.Reason)),
					Detail:           utils.Must(s.encrypt.DecryptString(report.Detail)),
					Images: slices.AppendSeq(
						make([]string, 0),
						utils.Map(report.Images, func(image string) string {
							return utils.Must(s.resourceService.SignURLWithDefaultDuration(image))
						}),
					),
					Status:        string(report.Status),
					AdminId:       report.AdminId.String,
					AdminUsername: utils.Must(s.encrypt.DecryptString(admin.Username)),
					AdminNote:     utils.Must(s.encrypt.DecryptString(report.AdminNote.String)),
					CreatedAt:     utils.FormatDMY(report.CreatedAt),
					CompletedAt:   utils.FormatDMY(report.UpdatedAt),
				}
			},
		),
	)

	booking := dto.Booking{}
	if report.Category == models.CATEGORY_BOOKING_RELATED {
		res, err := s.bookingStore.FindById(report.BookingId.String)
		if err != nil {
			return nil, err
		}

		service, err := s.serviceStore.FindById(res.ServiceId)
		if err != nil {
			return nil, err
		}

		booking = dto.Booking{
			Id: res.Id,
			Client: dto.User{
				Id:       reporter.Id,
				Name:     utils.Must(s.encrypt.DecryptString(reporter.Name)),
				Email:    utils.Must(s.encrypt.DecryptString(reporter.Email)),
				ImageURL: reporter.ImageUrl,
				Socials: slices.AppendSeq(
					make([]dto.Social, 0),
					utils.Map(reporter.Socials, func(social models.SocialModel) dto.Social {
						return dto.Social{
							Id:    social.Id,
							Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
							Title: utils.Must(s.encrypt.DecryptString(social.Title)),
							URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
						}
					}),
				),
				CreatedAt:    utils.FormatDate(reporter.CreatedAt),
				DateVerified: utils.FormatDate(reporter.VerifiedAt.String),
				IsRestricted: reporter.Restricted,
				IsBanned:     reporter.Banned,
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
			Quantity: res.Quantity,
			Cost:     res.Cost,
			Status:   string(res.Status),
			Extras: slices.AppendSeq(
				make([]dto.Extra, 0),
				utils.Map(res.Extras, func(x *models.ExtraModel) dto.Extra {
					return dto.Extra{
						Id:          x.Id,
						Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
						Description: utils.Must(s.encrypt.DecryptString(x.Description)),
						Price:       x.Price,
					}
				}),
			),
			CreatedAt:     utils.FormatDate(res.CreatedAt),
			ScheduleStart: utils.FormatDate(res.ScheduleStart.String),
			ScheduleEnd:   utils.FormatDate(res.ScheduleEnd.String),
			UpdatedAt:     utils.FormatDate(res.UpdatedAt),
			CancelReason:  utils.Try(s.encrypt.DecryptString(res.CancelReason.String)),
		}
	}

	data := &dto.UserReportDetail{
		Reporter: dto.User{
			Id:       reporter.Id,
			Name:     utils.Must(s.encrypt.DecryptString(reporter.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(reporter.Email)),
			ImageURL: reporter.ImageUrl,
			Address:  utils.Must(s.encrypt.DecryptString(reporter.Address.Address)),
			Phone:    utils.Must(s.encrypt.DecryptString(reporter.Phone)),
			Socials: slices.AppendSeq(
				make([]dto.Social, 0),
				utils.Map(reporter.Socials, func(social models.SocialModel) dto.Social {
					return dto.Social{
						Id:    social.Id,
						Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
						Title: utils.Must(s.encrypt.DecryptString(social.Title)),
						URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
					}
				}),
			),
			CreatedAt:    utils.FormatDate(reporter.CreatedAt),
			DateVerified: utils.FormatDate(reporter.VerifiedAt.String),
			IsRestricted: reporter.Restricted,
			IsBanned:     reporter.Banned,
		},
		Reported: dto.User{
			Id:       reported.Id,
			Name:     utils.Must(s.encrypt.DecryptString(reported.Name)),
			Email:    utils.Must(s.encrypt.DecryptString(reported.Email)),
			ImageURL: reported.ImageUrl,
			Address:  utils.Must(s.encrypt.DecryptString(reported.Address.Address)),
			Phone:    utils.Must(s.encrypt.DecryptString(reported.Phone)),
			Socials: slices.AppendSeq(
				make([]dto.Social, 0),
				utils.Map(reported.Socials, func(social models.SocialModel) dto.Social {
					return dto.Social{
						Id:    social.Id,
						Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
						Title: utils.Must(s.encrypt.DecryptString(social.Title)),
						URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
					}
				}),
			),
			CreatedAt:    utils.FormatDate(reported.CreatedAt),
			DateVerified: utils.FormatDate(reported.VerifiedAt.String),
			IsRestricted: reported.Restricted,
			IsBanned:     reported.Banned,
		},
		ReporterHistory: dto.ReporterHistory{
			ReportsFiled: len(reporterReportsFiled),
			FalseReports: len(
				slices.Collect(
					utils.Retain(reporterReportsFiled, func(r *models.UserReportModel) bool {
						return r.Status == models.REPORT_STATUS_DISMISSED
					}),
				),
			),
			CancelledBookings: len(reporterCancelledBooking),
			AccountCreatedAt:  utils.FormatDate(reporter.CreatedAt),
		},
		ReportedUserHistory: dto.ReportedUserHistory{
			Bookings: len(vendorBookings),
			CompletedBookings: len(
				slices.Collect(
					utils.Retain(vendorBookings, func(b *models.BookingModel) bool {
						return b.Status == models.BOOKING_STATUS_DONE
					}),
				),
			),
			RejectedBookings: len(
				slices.Collect(
					utils.Retain(vendorBookings, func(b *models.BookingModel) bool {
						return b.Status == models.BOOKING_STATUS_REJECTED
					}),
				),
			),
			ActiveBookings: len(
				slices.Collect(
					utils.Retain(vendorBookings, func(b *models.BookingModel) bool {
						return b.Status == models.BOOKING_STATUS_CONFIRMED
					}),
				),
			),
			PreviousReports:  reportHistory,
			AccountCreatedAt: utils.FormatDate(reported.CreatedAt),
			JoinedVendorAt:   utils.FormatDate(vendor.JoinedAt),
			Rating:           vendor.Rating,
		},
		Report: dto.ReportDetail{
			Id:               utils.FormatReportID(report.Id),
			ReportedByUserId: report.ReporterUserId,
			ReportedUserId:   report.ReportedUserId,
			Category:         string(report.Category),
			BookingId:        report.BookingId.String,
			Reason:           utils.Must(s.encrypt.DecryptString(report.Reason)),
			Detail:           utils.Must(s.encrypt.DecryptString(report.Detail)),
			Images: slices.AppendSeq(
				make([]string, 0),
				utils.Map(report.Images, func(image string) string {
					return utils.Must(s.resourceService.SignURLWithDefaultDuration(image))
				}),
			),
			Status:      string(report.Status),
			CreatedAt:   utils.FormatDate(report.CreatedAt),
			CompletedAt: utils.FormatDate(report.UpdatedAt),
		},
		Booking: booking,
	}

	return data, nil
}

// actions = "resolved" | "dismissed"
func (s *Service) ActOnReport(reportId, action, adminId, note string) error {
	allowedActions := []string{"resolved", "dismissed"}
	if !slices.Contains(allowedActions, action) {
		return errors.New("invalid_action")
	}

	parsedId, err := utils.ParseReportID(reportId)
	if err != nil {
		return err
	}

	// NOTE: action will serve as status. Refer to statuses, SHOULD MATCH
	encryptedNote := utils.Must(s.encrypt.EncryptString(note))
	if err := s.reportUserStore.Close(parsedId, action, adminId, encryptedNote); err != nil {
		return err
	}

	return nil
}
