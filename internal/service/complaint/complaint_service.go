package complaint_service

import (
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	notification_repo "nearbyassist/internal/repository/notification"
	report_user_repo "nearbyassist/internal/repository/report_user"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"strconv"
)

type Service struct {
	reportUserStore report_user_repo.ReportUserRepository
	userStore       user_repo.UserRepository
	bugReportStore  bug_report_repo.BugReportRepository
	notifStore      notification_repo.NotificationRepository
	ws              websocket.Socket
	fs              fs.FileStorage
	encrypt         core.Encryption
	jwt             core.Authenticator
}

func NewService(
	reportUserStore report_user_repo.ReportUserRepository,
	userStore user_repo.UserRepository,
	bugReportStore bug_report_repo.BugReportRepository,
	notifStore notification_repo.NotificationRepository,
	ws websocket.Socket,
	fs fs.FileStorage,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		reportUserStore: reportUserStore,
		userStore:       userStore,
		bugReportStore:  bugReportStore,
		notifStore:      notifStore,
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

func (s *Service) ReportUser(bearerToken string, req *request.ReportUserPayload, files []*multipart.FileHeader) (string, error) {
	reporterId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	reportData := &models.ReportedUserModel{
		ReportedBy: reporterId,
		UserId:     req.UserId,
		Reason:     utils.Must(s.encrypt.EncryptString(req.Reason)),
		Detail:     utils.Must(s.encrypt.EncryptString(req.Detail)),
		Images:     make([]string, 0),
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
			return "", err
		}

		reportData.Images = append(reportData.Images, url)
	}

	reportId, err := s.reportUserStore.Create(reportData)
	if err != nil {
		return "", err
	}

	return reportId, nil
}

func (s *Service) GetReportedUsers(limit, offset int) ([]*models.ReportedUserModel, error) {
	users, err := s.reportUserStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Reason = utils.Must(s.encrypt.DecryptString(user.Reason))
		user.Detail = utils.Must(s.encrypt.DecryptString(user.Detail))
	}

	return users, nil
}

func (s *Service) GetReportedUserDetail(reportId string) (*dto.UserReport, error) {
	report, err := s.reportUserStore.FindById(reportId)
	if err != nil {
		return nil, err
	}

	reporter, err := s.userStore.FindById(report.ReportedBy)
	if err != nil {
		return nil, err
	}

	reported, err := s.userStore.FindById(report.UserId)
	if err != nil {
		return nil, err
	}

	data := &dto.UserReport{
		Reporter: dto.User{
			Id:           reporter.Id,
			Name:         utils.Must(s.encrypt.DecryptString(reporter.Name)),
			Email:        utils.Must(s.encrypt.DecryptString(reporter.Email)),
			ImageURL:     reporter.ImageUrl,
			Address:      utils.Try(s.encrypt.DecryptString(reporter.Address.String)),
			Phone:        utils.Try(s.encrypt.DecryptString(reporter.Phone.String)),
			Socials:      reporter.Socials,
			CreatedAt:    utils.FormatDate(reporter.CreatedAt),
			DateVerified: utils.FormatDate(reporter.VerifiedAt),
			IsRestricted: reporter.Restricted,
			IsBanned:     reporter.Banned,
		},
		Reported: dto.User{
			Id:           reported.Id,
			Name:         utils.Must(s.encrypt.DecryptString(reported.Name)),
			Email:        utils.Must(s.encrypt.DecryptString(reported.Email)),
			ImageURL:     reported.ImageUrl,
			Address:      utils.Try(s.encrypt.DecryptString(reported.Address.String)),
			Phone:        utils.Try(s.encrypt.DecryptString(reported.Phone.String)),
			Socials:      reported.Socials,
			CreatedAt:    utils.FormatDate(reported.CreatedAt),
			DateVerified: utils.FormatDate(reported.VerifiedAt),
			IsRestricted: reported.Restricted,
			IsBanned:     reported.Banned,
		},
		Report: dto.Report{
			Id:               report.Id,
			ReportedByUserId: report.ReportedBy,
			ReportedUserId:   report.UserId,
			Reason:           utils.Must(s.encrypt.DecryptString(report.Reason)),
			Detail:           utils.Must(s.encrypt.DecryptString(report.Detail)),
			Images:           report.Images,
			Status:           "",
			CreatedAt:        utils.FormatDate(report.CreatedAt),
			CompletedAt:      utils.FormatDate(report.CompletedAt.String),
		},
	}

	return data, nil
}

func (s *Service) CloseUserReport(reportId, title, detail string) error {
	report, err := s.reportUserStore.FindById(reportId)
	if err != nil {
		return err
	}

	if err := s.reportUserStore.CloseReport(reportId); err != nil {
		return err
	}

	notificationHeading := "User report has been addressed!"
	notificationContent := "Your recent user report submission has been viewed and addressed!"

	notification := &models.NotificationModel{
		Recipient: report.ReportedBy,
		Type:      "generic",
		Title:     title,
		Content:   detail,
	}

	encryptedNotification := &models.NotificationModel{
		Recipient: report.ReportedBy,
		Type:      "generic",
		Title:     utils.Must(s.encrypt.EncryptString(notification.Title)),
		Content:   utils.Must(s.encrypt.EncryptString(notification.Content)),
	}

	if notifId, err := s.notifStore.Create(encryptedNotification); err != nil {
		return err
	} else {
		notification.Id = notifId
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewUrgentNotification(report.ReportedBy, notificationHeading, notificationContent); err != nil {
		fmt.Println(err.Error())
	}

	event := &websocket.EventModel{
		ReceiverId: report.ReportedBy,
		Type:       websocket.EVT_NOTIF,
		Payload:    notification,
	}

	s.ws.Send(event)

	return nil
}
