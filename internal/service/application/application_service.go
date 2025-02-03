package application_service

import (
	"encoding/base64"
	"mime/multipart"
	"nearbyassist/internal/models"
	application_repo "nearbyassist/internal/repository/application"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"
)

type Service struct {
	store   application_repo.ApplicationRepository
	fs      fs.FileStorage
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store application_repo.ApplicationRepository, fs fs.FileStorage, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{store: store, fs: fs, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateApplication(bearerToken, expertiseId string, files []*multipart.FileHeader) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	application := new(models.ApplicationModel)
	application.ApplicantId = userId
	application.ExpertiseId = expertiseId

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return "", err
		}

		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		switch file.Filename {
		case "policeClearance":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.POLICE_CLEARANCE_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return "", err
			}

			application.PoliceClearanceUrl = url

		case "supportingDocument":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.APPLICATION_PROOF_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return "", err
			}

			application.SupportingDocumentUrl = url

		default:
			return "", err

		}
	}

	applicationId, err := s.store.Create(application)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return "", err
		}

		return "", err
	}

	return applicationId, nil
}

func (s *Service) GetApplications() ([]*models.ApplicationModel, error) {
	applications, err := s.store.GetAll("pending")
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *Service) GetApplicationDetail(applicationId string) (*models.ApplicationModel, error) {
	application, err := s.store.FindApplication(applicationId)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(application.ApplicantName); err != nil {
		return nil, err
	} else {
		application.ApplicantName = decrypted
	}

	return application, nil
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
