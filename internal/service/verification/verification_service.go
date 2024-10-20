package verification_service

import (
	"mime/multipart"
	"nearbyassist/internal/models"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
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
