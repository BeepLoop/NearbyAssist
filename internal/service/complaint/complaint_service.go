package complaint_service

import (
	"mime/multipart"
	"nearbyassist/internal/models"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   bug_report_repo.BugReportRepository
	fs      fs.FileStorage
	encrypt auth.Encryption
}

func NewService(store bug_report_repo.BugReportRepository, fs fs.FileStorage, encrypt auth.Encryption) *Service {
	return &Service{store: store, fs: fs, encrypt: encrypt}
}

func (s *Service) CreateBugReport(req *request.BugReportPayload, files []*multipart.FileHeader) (string, error) {
	newComplaint := new(models.BugReportModel)
	newComplaint.Title = req.Title
	newComplaint.Detail = req.Detail

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
			Category: fs.SYS_COMPLAINT_DIR,
		}
		url, err := s.fs.SaveFile(fileData)
		if err != nil {
			return "", err
		}

		newComplaint.Images = append(newComplaint.Images, url)
	}

	if cipher, err := s.encrypt.EncryptString(req.Title); err != nil {
		return "", err
	} else {
		newComplaint.Title = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Detail); err != nil {
		return "", err
	} else {
		newComplaint.Detail = cipher
	}

	complaintId, err := s.store.Create(newComplaint)
	if err != nil {
		return "", err
	}

	return complaintId, nil
}

func (s *Service) GetBugReports(limit, offset int) ([]*models.BugReportModel, error) {
	complaints, err := s.store.GetAll(limit, offset)
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
