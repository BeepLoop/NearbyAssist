package e2ee_service

import (
	"errors"
	"nearbyassist/internal/models"
	e2ee_repo "nearbyassist/internal/repository/e2ee"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   e2ee_repo.E2EERepository
	encrypt core.Encryption
	jwt     core.Authenticator
}

func NewService(store e2ee_repo.E2EERepository, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{store: store, encrypt: encrypt, jwt: jwt}
}

func (s *Service) SaveKeys(bearerToken string, req *request.PEM) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	publicKey := new(models.PublicKeyModel)
	publicKey.Owner = userId
	publicKey.Pem = req.Public
	if _, err := s.store.NewPublicPem(publicKey); err != nil {
		return err
	}

	privateKey := new(models.PrivateKeyModel)
	privateKey.Owner = userId
	if encrypted, err := s.encrypt.EncryptString(req.Private); err != nil {
		return err
	} else {
		privateKey.Pem = encrypted
	}

	if _, err := s.store.NewPrivatePem(privateKey); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetKeys(bearerToken string) (map[string]string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	private, exists, err := s.store.GetPrivatePem(userId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("not found")
	}

	decrypted, err := s.encrypt.DecryptString(private.Pem)
	if err != nil {
		return nil, err
	}

	public, err := s.store.GetPublicPem(userId)
	if err != nil {
		return nil, err
	}

	keys := map[string]string{
		"publicKey":  public.Pem,
		"privateKey": decrypted,
	}

	return keys, nil
}

func (s *Service) GetPublicKey(userId string) (string, error) {
	publicKey, err := s.store.GetPublicPem(userId)
	if err != nil {
		return "", err
	}

	return publicKey.Pem, nil
}
