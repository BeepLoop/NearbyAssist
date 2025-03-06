package complaint_service

import (
	"mime/multipart"
	"nearbyassist/internal/models"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	report_user_repo "nearbyassist/internal/repository/report_user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
)

type Service struct {
	reportUserStore report_user_repo.ReportUserRepository
	bugReportStore  bug_report_repo.BugReportRepository
	fs              fs.FileStorage
	encrypt         auth.Encryption
}

func NewService(reportUserStore report_user_repo.ReportUserRepository, bugReportStore bug_report_repo.BugReportRepository, fs fs.FileStorage, encrypt auth.Encryption) *Service {
	return &Service{
		reportUserStore: reportUserStore,
		bugReportStore:  bugReportStore,
		fs:              fs,
		encrypt:         encrypt,
	}
}

func (s *Service) CreateBugReport(req *request.BugReportPayload, files []*multipart.FileHeader) (string, error) {
	reportData := &models.BugReportModel{
		Title:  req.Title,
		Detail: req.Detail,
		Images: make([]string, 0),
	}

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
		}

		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		fileData := fs.File{
			Data:     cipher,
			Category: fs.BUG_REPORT_DIR,
		}
		url, err := s.fs.SaveFile(fileData)
		if err != nil {
			return "", err
		}

		reportData.Images = append(reportData.Images, url)
	}

	if cipher, err := s.encrypt.EncryptString(req.Title); err != nil {
		return "", err
	} else {
		reportData.Title = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Detail); err != nil {
		return "", err
	} else {
		reportData.Detail = cipher
	}

	complaintId, err := s.bugReportStore.Create(reportData)
	if err != nil {
		return "", err
	}

	return complaintId, nil
}

func (s *Service) ReportUser(req *request.ReportUserPayload, files []*multipart.FileHeader) (string, error) {
	reportData := &models.ReportedUserModel{
		UserId: req.UserId,
		Title:  req.Title,
		Reason: req.Reason,
		Images: make([]string, 0),
	}

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
		}

		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		fileData := fs.File{
			Data:     cipher,
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

func (s *Service) GetBugReports(limit, offset int) ([]*models.BugReportModel, error) {
	complaints, err := s.bugReportStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, complaint := range complaints {
		if cipher, err := s.encrypt.DecryptString(complaint.Title); err != nil {
			return nil, err
		} else {
			complaint.Title = cipher
		}

		if cipher, err := s.encrypt.DecryptString(complaint.Detail); err != nil {
			return nil, err
		} else {
			complaint.Detail = cipher
		}
	}

	return complaints, nil
}
