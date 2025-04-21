package userauth_service

import (
	"errors"
	"mime/multipart"
	"nearbyassist/internal/models"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/repository/userauth"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
)

const (
	ERR_BANNED_USER  = "banned user"
	ERR_NOT_FOUND    = "not_found"
	ERR_EMAIL_EXISTS = "email already exists"
)

type Service struct {
	userStore         user_repo.UserRepository
	userAuthStore     userauth.Repository
	verificationStore verification_repo.VerificationRepository
	fs                fs.FileStorage
	encrypt           core.Encryption
	hash              core.Hash
	jwt               core.Authenticator
}

func NewService(
	userStore user_repo.UserRepository,
	userAuthStore userauth.Repository,
	verificationStore verification_repo.VerificationRepository,
	fs fs.FileStorage,
	encrypt core.Encryption,
	hash core.Hash,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:         userStore,
		userAuthStore:     userAuthStore,
		verificationStore: verificationStore,
		fs:                fs,
		encrypt:           encrypt,
		hash:              hash,
		jwt:               jwt,
	}
}

func (s *Service) Login(req *request.UserLoginPayload) (*response.LoginResponse, error) {
	emailHash, err := s.hash.Generate([]byte(req.Email))
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, errors.New(ERR_NOT_FOUND)
	}

	// Check if the user is banned
	if existingUser.Banned {
		return nil, errors.New(ERR_BANNED_USER)
	}
	existingUser.Name = utils.Must(s.encrypt.DecryptString(existingUser.Name))
	existingUser.Email = utils.Must(s.encrypt.DecryptString(existingUser.Email))
	existingUser.Address.Address = utils.Must(s.encrypt.DecryptString(existingUser.Address.Address))
	existingUser.Phone = utils.Must(s.encrypt.DecryptString(existingUser.Phone))

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
	if err := s.userAuthStore.Login(session); err != nil {
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
			Address:      existingUser.Address.Address,
			Phone:        existingUser.Phone,
			Latitude:     existingUser.Address.Latitude,
			Longitude:    existingUser.Address.Longitude,
			Expertises:   vendorExpertises,
			IsRestricted: existingUser.Restricted,
		},
	}

	return response, nil
}

func (s *Service) Register(req *request.UserRegisterPayload, files []*multipart.FileHeader) (*response.LoginResponse, error) {
	emailHash, err := s.hash.Generate([]byte(req.Email))
	if err != nil {
		return nil, err
	}

	existing, _ := s.userStore.FindByEmailHash(emailHash)
	if existing != nil {
		return nil, errors.New(ERR_EMAIL_EXISTS)
	}

	user := &models.UserModel{
		Name:      utils.Must(s.encrypt.EncryptString(req.Name)),
		Email:     utils.Must(s.encrypt.EncryptString(req.Email)),
		ImageUrl:  req.ImageURL,
		EmailHash: emailHash,
		Phone:     utils.Must(s.encrypt.EncryptString(req.Phone)),
		Identification: models.IdentificationModel{
			Type:            req.IDType,
			ReferenceNumber: utils.Must(s.encrypt.EncryptString(req.ReferenceNumber)),
		},
		Address: models.AddressModel{
			Address:   utils.Must(s.encrypt.EncryptString(req.Address)),
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
	}

	for _, file := range files {
		// Read bytes
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return nil, err
		}

		// Encrypt the file
		cipher, err := s.encrypt.EncryptFile(bytes)
		if err != nil {
			return nil, err
		}

		switch file.Filename {
		case "frontId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_FRONT,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return nil, err
			} else {
				user.Identification.FrontImageUrl = url
			}

		case "backId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_BACK,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return nil, err
			} else {
				user.Identification.BackImageUrl = url
			}

		case "face":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.FACE,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return nil, err
			} else {
				user.Identification.SelfieImageUrl = url
			}

		default:
			return nil, err
		}
	}

	userId, err := s.userStore.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if _, err := s.verificationStore.CreateLink(userId); err != nil {
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
	if err := s.userAuthStore.Login(session); err != nil {
		return nil, err
	}

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:           user.Id,
			Name:         req.Name,
			Email:        req.Email,
			ImageUrl:     req.ImageURL,
			IsVerified:   false,
			IsVendor:     false,
			Address:      req.Address,
			Phone:        req.Phone,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
			Expertises:   make([]response.Expertise, 0),
			IsRestricted: false,
		},
	}

	return response, nil
}

func (s *Service) Refresh(bearerToken, refreshToken string) (string, error) {
	// Check if refreshToken exists
	if _, err := s.userAuthStore.FindSessionByToken(refreshToken); err != nil {
		return "", err
	}

	// Check if refreshToken is blacklisted
	if err := s.userAuthStore.IsRefreshTokenBlacklisted(refreshToken); err == nil {
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
	if _, err := s.userAuthStore.FindSessionByToken(refreshToken); err != nil {
		return err
	}

	if err := s.userAuthStore.Logout(refreshToken); err != nil {
		return err
	}

	return nil
}
