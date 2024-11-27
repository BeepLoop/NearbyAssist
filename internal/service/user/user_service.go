package user_service

import (
	"database/sql"
	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
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

func (s *Service) Login(req *request.UserLoginPayload) (*response.LoginResponse, error) {
	emailHash, err := s.hash.Generate([]byte(req.Email))
	if err != nil {
		return nil, err
	}

	existingUser, err := s.store.FindByEmailHash(emailHash)
	if err != nil {
		// If user is not found, continue to registration
		return s.Register(req, emailHash)
	}

	if existingUser.Address.Valid {
		if plain, err := s.encrypt.DecryptString(existingUser.Address.String); err != nil {
			return nil, err
		} else {
			existingUser.Address = sql.NullString{String: plain, Valid: true}
		}
	}

	if existingUser.Latitude.Valid == false {
		existingUser.Latitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	if existingUser.Longitude.Valid == false {
		existingUser.Longitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	isVendor, err := s.store.IsVendor(existingUser.Id)
	if err != nil {
		return nil, err
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

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:         existingUser.Id,
			Name:       req.Name,
			Email:      req.Email,
			ImageUrl:   req.Image,
			IsVerified: existingUser.Verified,
			IsVendor:   isVendor,
			Address:    existingUser.Address.String,
			Latitude:   existingUser.Latitude.Float64,
			Longitude:  existingUser.Longitude.Float64,
		},
	}

	return response, nil
}

func (s *Service) Register(req *request.UserLoginPayload, emailHash string) (*response.LoginResponse, error) {
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

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:         newUser.Id,
			Name:       req.Name,
			Email:      req.Email,
			ImageUrl:   req.Image,
			IsVerified: false,
			IsVendor:   false,
		},
	}

	return response, nil
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

func (s *Service) GetUser(bearerToken string) (*response.DetailedUser, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		return nil, err
	}

	isVendor, err := s.store.IsVendor(user.Id)
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

	if user.Address.Valid {
		if plain, err := s.encrypt.DecryptString(user.Address.String); err != nil {
			return nil, err
		} else {
			user.Address = sql.NullString{String: plain, Valid: true}
		}
	}

	if user.Latitude.Valid == false {
		user.Latitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	if user.Longitude.Valid == false {
		user.Longitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	response := &response.DetailedUser{
		Id:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		IsVerified: user.Verified,
		IsVendor:   isVendor,
		Address:    user.Address.String,
		Latitude:   user.Latitude.Float64,
		Longitude:  user.Longitude.Float64,
	}

	return response, nil
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
