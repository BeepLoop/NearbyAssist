package resource_service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DEFAULT_DURATION = time.Minute * 1
)

type Service struct {
	fs      fs.FileStorage
	encrypt auth.Encryption
	hash    auth.Hash
}

func NewService(fs fs.FileStorage, encrypt auth.Encryption, hash auth.Hash) *Service {
	return &Service{
		fs:      fs,
		encrypt: encrypt,
		hash:    hash,
	}
}

func (s *Service) SignURL(imagePath string, duration time.Duration) (string, error) {
	expiry := time.Now().Add(duration).Unix()
	data := fmt.Sprintf("%s:%d", imagePath, expiry)

	signature, err := s.hash.Generate([]byte(data))
	if err != nil {
		return "", err
	}

	signedURL := fmt.Sprintf("image?path=%s&expiry=%d&signature=%s", url.QueryEscape(imagePath), expiry, signature)

	return signedURL, nil
}

// Returns raw bytes
func (s *Service) GetRawFile(path string) ([]byte, error) {
	b, err := s.fs.GetFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetPrivateFile(path, signature, expiry string) ([]byte, error) {
	// Validate expiry
	expireTime, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil || time.Now().Unix() > expireTime {
		return nil, errors.New("unauthorized access")
	}

	// Validate signature
	data := fmt.Sprintf("%s:%s", path, expiry)
	expectedSignature, err := s.hash.Generate([]byte(data))
	if err != nil {
		return nil, err
	}

	if signature != expectedSignature {
		return nil, errors.New("unauthorized access")
	}

	// Retrieve file
	file, err := s.fs.GetFile(path)

	decrypted, err := s.encrypt.DecryptFile(file)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// Returns decrypted file in base64 form
func (s *Service) GetBase64File(path string) (string, error) {
	file, err := s.fs.GetFile(path)

	decrypted, err := s.encrypt.DecryptFile(file)
	if err != nil {
		return "", err
	}

	return s.ImageToBase64(decrypted)
}

func (s *Service) ImageToBase64(file []byte) (string, error) {
	base64Img := base64.StdEncoding.EncodeToString(file)

	mime := http.DetectContentType(file)

	base64Img = "data:" + mime + ";base64," + base64Img

	return base64Img, nil
}
