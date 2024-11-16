package verification_service

import (
	"encoding/base64"
	"mime/multipart"
	"nearbyassist/internal/models"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
	"net/http"
)

type Service struct {
	store   verification_repo.VerificationRepository
	fs      fs.FileStorage
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store verification_repo.VerificationRepository, fs fs.FileStorage, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{store: store, fs: fs, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateVerificationRequest(name, address, idType, idNumber, bearerToken string, files []*multipart.FileHeader) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	encryptedName, err := s.encrypt.EncryptString(name)
	if err != nil {
		return "", err
	}

	encryptedAddress, err := s.encrypt.EncryptString(address)
	if err != nil {
		return "", err
	}

	encryptedIdNumber, err := s.encrypt.EncryptString(idNumber)
	if err != nil {
		return "", err
	}

	req := new(models.IdentityVerificationModel)
	req.UserId = userId
	req.Name = encryptedName
	req.Address = encryptedAddress
	req.IdType = idType
	req.IdNumber = encryptedIdNumber

	for _, file := range files {
		// Read bytes
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return "", err
		}

		// Encrypt the file
		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return "", err
		}

		switch file.Filename {
		case "frontId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_FRONT,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.FrontIdImageUrl = url
			}

		case "backId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_BACK,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.BackIdImageUrl = url
			}

		case "face":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.FACE,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return "", err
			} else {
				req.FaceImageUrl = url
			}

		default:
			return "", err
		}
	}

	verificationId, err := s.store.Create(req)
	if err != nil {
		return "", err
	}

	return verificationId, nil
}

func (s *Service) GetRequest(id string) (*models.IdentityVerificationModel, error) {
	request, err := s.store.FindById(id)
	if err != nil {
		return nil, err
	}

	decryptedName, err := s.encrypt.DecryptString(request.Name)
	if err != nil {
		return nil, err
	}
	request.Name = decryptedName

	decryptedAddress, err := s.encrypt.DecryptString(request.Address)
	if err != nil {
		return nil, err
	}
	request.Address = decryptedAddress

	decryptedIdNumber, err := s.encrypt.DecryptString(request.IdNumber)
	if err != nil {
		return nil, err
	}
	request.IdNumber = decryptedIdNumber

	return request, nil
}

func (s *Service) GetIdentityVerificationRequests() ([]*models.IdentityVerificationModel, error) {
	requests, err := s.store.GetAll("pending")
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		decryptedName, err := s.encrypt.DecryptString(request.Name)
		if err != nil {
			return nil, err
		}
		request.Name = decryptedName

		decryptedAddress, err := s.encrypt.DecryptString(request.Address)
		if err != nil {
			return nil, err
		}
		request.Address = decryptedAddress

		decryptedIdNumber, err := s.encrypt.DecryptString(request.IdNumber)
		if err != nil {
			return nil, err
		}
		request.IdNumber = decryptedIdNumber
	}

	return requests, nil
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

func (s *Service) AcceptRequest(id string) error {
	return s.store.AcceptRequest(id)
}

func (s *Service) RejectRequest(id string) error {
	return s.store.RejectRequest(id)
}
