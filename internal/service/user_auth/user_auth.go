package userauth_service

import (
	"database/sql"
	"errors"
	"nearbyassist/internal/models"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

const (
	ERR_BANNED_USER = "banned user"
)

type Service struct {
	userStore user_repo.UserRepository
	encrypt   core.Encryption
	hash      core.Hash
	jwt       core.Authenticator
}

func NewService(userStore user_repo.UserRepository, encrypt core.Encryption, hash core.Hash, jwt core.Authenticator) *Service {
	return &Service{
		userStore: userStore,
		encrypt:   encrypt,
		hash:      hash,
		jwt:       jwt,
	}
}

func (s *Service) ThirdPartyLogin(req *request.UserLoginPayload) (*response.LoginResponse, error) {
	emailHash, err := s.hash.Generate([]byte(req.Email))
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userStore.FindByEmailHash(emailHash)
	if err != nil {
		// If user is not found, continue to registration
		return s.ThirdPartyRegister(req, emailHash)
	}

	// Check if the user is banned
	if existingUser.Banned {
		return nil, errors.New(ERR_BANNED_USER)
	}

	if decrypted, err := s.encrypt.DecryptString(existingUser.Name); err != nil {
		return nil, err
	} else {
		existingUser.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(existingUser.Email); err != nil {
		return nil, err
	} else {
		existingUser.Email = decrypted
	}

	if existingUser.Address.Valid {
		if plain, err := s.encrypt.DecryptString(existingUser.Address.String); err != nil {
			return nil, err
		} else {
			existingUser.Address = sql.NullString{String: plain, Valid: true}
		}
	}

	if existingUser.Phone.Valid {
		if plain, err := s.encrypt.DecryptString(existingUser.Phone.String); err != nil {
			return nil, err
		} else {
			existingUser.Phone = sql.NullString{String: plain, Valid: true}
		}
	}

	if existingUser.Latitude.Valid == false {
		existingUser.Latitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	if existingUser.Longitude.Valid == false {
		existingUser.Longitude = sql.NullFloat64{Float64: 0.0, Valid: true}
	}

	isVendor, err := s.userStore.IsVendor(existingUser.Id)
	if err != nil {
		return nil, err
	}

	vendorExpertises := make([]response.Expertise, 0)
	if isVendor {
		expertises, err := s.userStore.GetExpertise(existingUser.Id)
		if err != nil {
			return nil, err
		}

		for _, expertise := range expertises {
			expertiseTags := make([]response.Tag, 0)

			for _, tag := range expertise.Tags {
				expertiseTags = append(expertiseTags, response.Tag{
					Id:    tag.Id,
					Title: tag.Title,
				})
			}

			vendorExpertises = append(vendorExpertises, response.Expertise{
				Id:    expertise.Id,
				Title: expertise.Title,
				Tags:  expertiseTags,
			})
		}
	}

	accessToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: existingUser.Id,
		Name:   existingUser.Name,
		Email:  existingUser.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := models.NewSessionModel(refreshToken)
	if err := s.userStore.Login(session); err != nil {
		return nil, err
	}

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:           existingUser.Id,
			Name:         existingUser.Name,
			Email:        existingUser.Email,
			ImageUrl:     existingUser.ImageUrl,
			IsVerified:   existingUser.Verified,
			IsVendor:     isVendor,
			Address:      existingUser.Address.String,
			Phone:        existingUser.Phone.String,
			Latitude:     existingUser.Latitude.Float64,
			Longitude:    existingUser.Longitude.Float64,
			Expertises:   vendorExpertises,
			IsRestricted: existingUser.Restricted,
		},
	}

	return response, nil
}

func (s *Service) ThirdPartyRegister(req *request.UserLoginPayload, emailHash string) (*response.LoginResponse, error) {
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

	userId, err := s.userStore.CreateUser(newUser)
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
	if err := s.userStore.Login(session); err != nil {
		return nil, err
	}

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:           newUser.Id,
			Name:         req.Name,
			Email:        req.Email,
			ImageUrl:     req.Image,
			IsVerified:   false,
			IsVendor:     false,
			IsRestricted: false,
		},
	}

	return response, nil
}

func (s *Service) Refresh(bearerToken, refreshToken string) (string, error) {
	// Check if refreshToken exists
	if _, err := s.userStore.FindSessionByToken(refreshToken); err != nil {
		return "", err
	}

	// Check if refreshToken is blacklisted
	if err := s.userStore.IsRefreshTokenBlacklisted(refreshToken); err == nil {
		return "", err
	}

	// Generate new accessToken
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	user, err := s.userStore.FindById(userId)
	if err != nil {
		return "", err
	}

	if user.Banned {
		return "", errors.New(ERR_BANNED_USER)
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
	if _, err := s.userStore.FindSessionByToken(refreshToken); err != nil {
		return err
	}

	if err := s.userStore.Logout(refreshToken); err != nil {
		return err
	}

	return nil
}
