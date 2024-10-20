package user_service

import (
	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   repository.UserRepository
	encrypt auth.Encryption
	hash    auth.Hash
	jwt     auth.Authenticator
}

func NewService(store repository.UserRepository, encrypt auth.Encryption, hash auth.Hash, jwt auth.Authenticator) *Service {
	return &Service{store: store, encrypt: encrypt, hash: hash, jwt: jwt}
}

func (s *Service) Login(req *request.UserLoginPayload) (map[string]interface{}, error) {
	emailHash, err := s.hash.Generate([]byte(req.Email))
	if err != nil {
		return nil, err
	}

	existingUser, err := s.store.FindByEmailHash(emailHash)
	if err != nil {
		// If user is not found, continue to registration
		return s.Register(req, emailHash)
	}

	accessToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: existingUser.Id,
		Name:   req.Name,
		Email:  req.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := models.NewSessionModel(refreshToken)
	if err := s.store.Login(session); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"user": models.UserModel{
			Model:    models.Model{Id: existingUser.Id},
			Name:     req.Name,
			Email:    req.Email,
			ImageUrl: req.Image,
			Verified: existingUser.Verified,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return data, nil
}

func (s *Service) Register(req *request.UserLoginPayload, emailHash string) (map[string]interface{}, error) {
	newUser := new(models.UserModel)
	newUser.EmailHash = emailHash
	newUser.ImageUrl = req.Image

	if cipher, err := s.encrypt.EncryptString(req.Name); err != nil {
		return nil, err
	} else {
		newUser.Name = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Email); err != nil {
		return nil, err
	} else {
		newUser.Email = cipher
	}

	userId, err := s.store.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: userId,
		Name:   req.Name,
		Email:  req.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := models.NewSessionModel(refreshToken)
	if err := s.store.Login(session); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"user": models.UserModel{
			Model:    models.Model{Id: newUser.Id},
			Name:     req.Name,
			Email:    req.Email,
			ImageUrl: req.Image,
			Verified: newUser.Verified,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return data, nil
}

func (s *Service) Refresh(bearerToken, refreshToken string) (string, error) {
	// Check if refreshToken exists
	if _, err := s.store.FindSessionByToken(refreshToken); err != nil {
		return "", err
	}

	// Check if refreshToken is blacklisted
	if err := s.store.IsRefreshTokenBlacklisted(refreshToken); err == nil {
		return "", err
	}

	// Generate new accessToken
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		return "", err
	}

	if plain, err := s.encrypt.DecryptString(user.Name); err != nil {
		return "", err
	} else {
		user.Name = plain
	}

	if plain, err := s.encrypt.DecryptString(user.Email); err != nil {
		return "", err
	} else {
		user.Email = plain
	}

	newToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: user.Id,
		Name:   user.Name,
		Email:  user.Email,
	})
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *Service) Logout(refreshToken string) error {
	if _, err := s.store.FindSessionByToken(refreshToken); err != nil {
		return err
	}

	if err := s.store.Logout(refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUser(bearerToken string) (*models.UserModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		return nil, err
	}

	if plain, err := s.encrypt.DecryptString(user.Name); err != nil {
		return nil, err
	} else {
		user.Name = plain
	}

	if plain, err := s.encrypt.DecryptString(user.Email); err != nil {
		return nil, err
	} else {
		user.Email = plain
	}

	return user, nil
}

func (s *Service) IsVerified(bearerToken string) (bool, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return false, err
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		return false, err
	}

	return user.Verified, nil
}
