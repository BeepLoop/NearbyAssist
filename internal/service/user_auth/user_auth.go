package userauth_service

import (
	"errors"
	"nearbyassist/internal/models"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/repository/userauth"
	vendor_repo "nearbyassist/internal/repository/vendor"
	verification_repo "nearbyassist/internal/repository/verification"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	ERR_BANNED_USER  = "banned user"
	ERR_NOT_FOUND    = "not_found"
	ERR_EMAIL_EXISTS = "email already exists"
)

type Service struct {
	userStore         user_repo.UserRepository
	vendorStore       vendor_repo.VendorRepository
	userAuthStore     userauth.Repository
	verificationStore verification_repo.VerificationRepository
	fs                fs.FileStorage
	encrypt           core.Encryption
	hash              core.Hash
	jwt               core.Authenticator
}

func NewService(
	userStore user_repo.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	userAuthStore userauth.Repository,
	verificationStore verification_repo.VerificationRepository,
	fs fs.FileStorage,
	encrypt core.Encryption,
	hash core.Hash,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:         userStore,
		vendorStore:       vendorStore,
		userAuthStore:     userAuthStore,
		verificationStore: verificationStore,
		fs:                fs,
		encrypt:           encrypt,
		hash:              hash,
		jwt:               jwt,
	}
}

func (s *Service) Login(req *request.UserLoginPayload) (*response.LoginResponse, error) {
	emailHash := utils.Must(s.hash.Generate([]byte(req.Email)))
	user, err := s.userStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, errors.New(ERR_NOT_FOUND)
	}

	// Check if the user is banned
	if user.Banned {
		return nil, errors.New(ERR_BANNED_USER)
	}

	isVendor, err := s.userStore.IsVendor(user.Id)
	if err != nil {
		return nil, err
	}

	dailyBookingLimit := 0
	vendorExpertises := make([]response.Expertise, 0)
	if isVendor {
		if vendor, err := s.vendorStore.FindById(user.Id); err != nil {
			return nil, err
		} else {
			dailyBookingLimit = vendor.DBL
		}

		expertises, err := s.userStore.GetExpertise(user.Id)
		if err != nil {
			return nil, err
		}

		vendorExpertises = slices.AppendSeq(
			make([]response.Expertise, 0),
			utils.Map(expertises, func(expertise *models.ExpertiseModel) response.Expertise {
				return response.Expertise{
					Id:    expertise.Id,
					Title: expertise.Title,
					Tags: slices.AppendSeq(
						make([]response.Tag, 0),
						utils.Map(expertise.Tags, func(tag *models.TagModel) response.Tag {
							return response.Tag{
								Id:    tag.Id,
								Title: tag.Title,
							}
						}),
					),
				}
			}),
		)
	}

	accessToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: user.Id,
		Name:   user.Name,
		Email:  user.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken := utils.Must(s.jwt.GenerateRefreshToken())
	session := models.NewSessionModel(refreshToken)
	if err := s.userAuthStore.Login(session); err != nil {
		return nil, err
	}

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:         user.Id,
			Name:       utils.Must(s.encrypt.DecryptString(user.Name)),
			Email:      utils.Must(s.encrypt.DecryptString(user.Email)),
			ImageUrl:   user.ImageUrl,
			IsVerified: user.Verified,
			IsVendor:   isVendor,
			Address:    utils.Must(s.encrypt.DecryptString(user.Address.Address)),
			Phone:      utils.Must(s.encrypt.DecryptString(user.Phone)),
			Latitude:   user.Address.Latitude,
			Longitude:  user.Address.Longitude,
			Expertises: vendorExpertises,
			Socials: slices.AppendSeq(
				make([]response.Social, 0),
				utils.Map(user.Socials, func(social models.SocialModel) response.Social {
					return response.Social{
						Id:    social.Id,
						Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
						Title: utils.Must(s.encrypt.DecryptString(social.Title)),
						URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
					}
				}),
			),
			IsRestricted:           user.Restricted,
			DBL:                    dailyBookingLimit,
			HasPendingVerification: user.HasPendingVerification,
			HasPendingApplication:  user.HasPendingApplication,
		},
	}

	return response, nil
}

func (s *Service) Register(req *request.UserRegisterPayload) (*response.LoginResponse, error) {
	emailHash := utils.Must(s.hash.Generate([]byte(req.Email)))
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
		Address: models.AddressModel{
			Address:   utils.Must(s.encrypt.EncryptString(req.Address)),
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
	}

	userId, err := s.userStore.CreateUser(user)
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

	refreshToken := utils.Must(s.jwt.GenerateRefreshToken())
	session := models.NewSessionModel(refreshToken)

	if err := s.userAuthStore.Login(session); err != nil {
		return nil, err
	}

	response := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: response.DetailedUser{
			Id:                     user.Id,
			Name:                   req.Name,
			Email:                  req.Email,
			ImageUrl:               req.ImageURL,
			IsVerified:             false,
			IsVendor:               false,
			Address:                req.Address,
			Phone:                  req.Phone,
			Latitude:               req.Latitude,
			Longitude:              req.Longitude,
			Expertises:             make([]response.Expertise, 0),
			Socials:                make([]response.Social, 0),
			IsRestricted:           false,
			DBL:                    0,
			HasPendingVerification: false,
			HasPendingApplication:  false,
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

	newToken, err := s.jwt.GenerateAccessToken(models.JWTClaims{
		UserId: user.Id,
		Name:   utils.Must(s.encrypt.DecryptString(user.Name)),
		Email:  utils.Must(s.encrypt.DecryptString(user.Email)),
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
