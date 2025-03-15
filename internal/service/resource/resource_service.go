package resource_service

import (
	"encoding/base64"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"net/http"
)

type Service struct {
	fs      fs.FileStorage
	encrypt auth.Encryption
}

func NewService(fs fs.FileStorage, encrypt auth.Encryption) *Service {
	return &Service{
		fs:      fs,
		encrypt: encrypt,
	}
}

func (s *Service) GetFile(path string) ([]byte, error) {
	b, err := s.fs.GetFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetBase64File(path string) (string, error) {
	file, err := s.fs.GetFile(path)

	decrypted, err := s.encrypt.DecryptFile(file)
	if err != nil {
		return "", err
	}

	base64Img := base64.StdEncoding.EncodeToString(decrypted)

	mime := http.DetectContentType(decrypted)

	base64Img = "data:" + mime + ";base64," + base64Img

	return base64Img, nil
}
