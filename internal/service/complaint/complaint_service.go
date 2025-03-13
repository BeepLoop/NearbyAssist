package complaint_service

import (
	"encoding/base64"
	"mime/multipart"
	"nearbyassist/internal/models"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	report_user_repo "nearbyassist/internal/repository/report_user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"
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

	if cipher, err := s.encrypt.EncryptString(req.Title); err != nil {
		return err
	} else {
		reportData.Title = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Detail); err != nil {
		return err
	} else {
		reportData.Detail = cipher
	}

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

func (s *Service) CompleteBug(bugId string) error {
	id, err := strconv.Atoi(bugId)
	if err != nil {
		return err
	}

	return s.bugReportStore.CompleteBug(id)
}

func (s *Service) ReportUser(req *request.ReportUserPayload, files []*multipart.FileHeader) (string, error) {
	reason, err := s.encrypt.EncryptString(req.Reason)
	if err != nil {
		return "", err
	}

	detail, err := s.encrypt.EncryptString(req.Detail)
	if err != nil {
		return "", err
	}

	reportData := &models.ReportedUserModel{
		UserId: req.UserId,
		Reason: reason,
		Detail: detail,
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

func (s *Service) GetReportedUsers(limit, offset int) ([]*models.ReportedUserModel, error) {
	users, err := s.reportUserStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if decrypted, err := s.encrypt.DecryptString(user.Reason); err != nil {
			return nil, err
		} else {
			user.Reason = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(user.Detail); err != nil {
			return nil, err
		} else {
			user.Detail = decrypted
		}
	}

	return users, nil
}

func (s *Service) GetReportedUserDetail(reportId string) (*models.ReportedUserModel, error) {
	report, err := s.reportUserStore.FindById(reportId)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(report.Name); err != nil {
		return nil, err
	} else {
		report.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(report.Reason); err != nil {
		return nil, err
	} else {
		report.Reason = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(report.Detail); err != nil {
		return nil, err
	} else {
		report.Detail = decrypted
	}

	return report, nil
}

func (s *Service) GetFile(path string) (string, error) {
	file, err := s.fs.GetFile(path)
	if err != nil {
		return "", err
	}

	decrypted, err := s.encrypt.DecryptFile(file)
	if err != nil {
		return "", err
	}

	base64Img := base64.StdEncoding.EncodeToString(decrypted)

	mime := http.DetectContentType(decrypted)

	base64Img = "data:" + mime + ";base64," + base64Img

	return base64Img, nil
}
