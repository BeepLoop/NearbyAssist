package qr_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"nearbyassist/internal/request"
)

type Service struct {
	key []byte
}

func NewService(key []byte) *Service {
	return &Service{key: key}
}

func (s *Service) SignData(input *request.QRSignatureInput) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, s.key)

	if _, err := h.Write(data); err != nil {
		return "", err
	}

	signed := hex.EncodeToString(h.Sum(nil))

	return signed, nil
}

func (s *Service) VerifySignature(input *request.QRSignatureVerifyInput) bool {
	receivedSignature := input.Signature

	data := &request.QRSignatureInput{
		ClientID:  input.ClientID,
		VendorID:  input.VendorID,
		BookingID: input.BookingID,
	}

	expectedSignature, err := s.SignData(data)
	if err != nil {
		return false
	}

	return hmac.Equal([]byte(expectedSignature), []byte(receivedSignature))
}
