package user_service

import (
	"errors"
	"mime/multipart"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	application_repo "nearbyassist/internal/repository/application"
	supportingimage_repo "nearbyassist/internal/repository/supporting_image"
	repository "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/utils"
	"slices"
	"strconv"
)

const (
	ERR_FORBIDDEN             = "forbidden action"
	ERR_DUPLICATE_EXPERTISE   = "already have expertise"
	ERR_DUPLICATE_APPLICATION = "already have pending application for the expertise"
	ERR_INVALID_DBL           = "invalid dbl value"
)

type Service struct {
	userStore            repository.UserRepository
	vendorStore          vendor_repo.VendorRepository
	applicationStore     application_repo.ApplicationRepository
	supportingImageStore supportingimage_repo.Repository
	resourceService      *resource_service.Service
	fs                   fs.FileStorage
	encrypt              core.Encryption
	hash                 core.Hash
	jwt                  core.Authenticator
}

func NewService(
	userStore repository.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	applicationStore application_repo.ApplicationRepository,
	supportingImageStore supportingimage_repo.Repository,
	resourceService *resource_service.Service,
	fs fs.FileStorage,
	encrypt core.Encryption,
	hash core.Hash,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:            userStore,
		vendorStore:          vendorStore,
		applicationStore:     applicationStore,
		supportingImageStore: supportingImageStore,
		resourceService:      resourceService,
		fs:                   fs,
		encrypt:              encrypt,
		hash:                 hash,
		jwt:                  jwt,
	}
}

func (s *Service) GetAll(limit, offset int) ([]*models.UserModel, error) {
	accounts, err := s.userStore.GetAllUserAccounts(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if decrypted, err := s.encrypt.DecryptString(account.Name); err != nil {
			return nil, err
		} else {
			account.Name = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(account.Email); err != nil {
			return nil, err
		} else {
			account.Email = decrypted
		}
	}

	return accounts, nil
}

func (s *Service) GetAllBasicUsers(limit, offset int) ([]dto.User, error) {
	accounts, err := s.userStore.GetBasicUserAccounts(limit, offset)
	if err != nil {
		return nil, err
	}

	data := slices.AppendSeq(
		make([]dto.User, 0),
		utils.Map(accounts, func(account *models.UserModel) dto.User {
			identification := dto.Identification{}
			if account.HasSubmittedIdentification {
				identification = dto.Identification{
					Type:          account.Identification.Type,
					IdNumber:      utils.Must(s.encrypt.DecryptString(account.Identification.ReferenceNumber)),
					FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(account.Identification.FrontImageUrl)),
					BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(account.Identification.BackImageUrl)),
				}
			}

			return dto.User{
				Id:       account.Id,
				Name:     utils.Must(s.encrypt.DecryptString(account.Name)),
				Email:    utils.Must(s.encrypt.DecryptString(account.Email)),
				ImageURL: account.ImageUrl,
				Address:  utils.Must(s.encrypt.DecryptString(account.Address.Address)),
				Phone:    utils.Must(s.encrypt.DecryptString(account.Phone)),
				Socials: slices.AppendSeq(
					make([]dto.Social, 0),
					utils.Map(account.Socials, func(social models.SocialModel) dto.Social {
						return dto.Social{
							Id:    social.Id,
							Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
							Title: utils.Must(s.encrypt.DecryptString(social.Title)),
							URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
						}
					}),
				),
				Identification:             identification,
				CreatedAt:                  utils.FormatDate(account.CreatedAt),
				DateVerified:               utils.FormatDate(account.VerifiedAt.String),
				IsRestricted:               account.Restricted,
				IsBanned:                   account.Banned,
				HasSubmittedIdentification: account.HasSubmittedIdentification,
			}
		}),
	)

	return data, nil
}

func (s *Service) FindById(id string) (*dto.User, error) {
	user, err := s.userStore.FindById(id)
	if err != nil {
		return nil, err
	}

	identification := dto.Identification{}
	if user.HasSubmittedIdentification {
		identification = dto.Identification{
			Type:          user.Identification.Type,
			IdNumber:      utils.Must(s.encrypt.DecryptString(user.Identification.ReferenceNumber)),
			FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.FrontImageUrl)),
			BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.BackImageUrl)),
		}
	}

	data := &dto.User{
		Id:       user.Id,
		Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
		Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
		ImageURL: user.ImageUrl,
		Address:  utils.Must(s.encrypt.DecryptString(user.Address.Address)),
		Phone:    utils.Must(s.encrypt.DecryptString(user.Phone)),
		Socials: slices.AppendSeq(
			make([]dto.Social, 0),
			utils.Map(user.Socials, func(social models.SocialModel) dto.Social {
				return dto.Social{
					Id:    social.Id,
					Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
					Title: utils.Must(s.encrypt.DecryptString(social.Title)),
					URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
				}
			}),
		),
		Identification:             identification,
		CreatedAt:                  utils.FormatDate(user.CreatedAt),
		DateVerified:               utils.FormatDate(user.VerifiedAt.String),
		IsRestricted:               user.Restricted,
		IsBanned:                   user.Banned,
		HasSubmittedIdentification: user.HasSubmittedIdentification,
	}

	return data, nil
}

func (s *Service) FindByEmail(email string) (*dto.User, error) {
	emailHash := utils.Must(s.hash.Generate([]byte(email)))
	user, err := s.userStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, err
	}

	identification := dto.Identification{}
	if user.HasSubmittedIdentification {
		identification = dto.Identification{
			Type:          user.Identification.Type,
			IdNumber:      utils.Must(s.encrypt.DecryptString(user.Identification.ReferenceNumber)),
			FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.FrontImageUrl)),
			BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.BackImageUrl)),
		}
	}

	data := &dto.User{
		Id:       user.Id,
		Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
		Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
		ImageURL: user.ImageUrl,
		Address:  utils.Must(s.encrypt.DecryptString(user.Address.Address)),
		Phone:    utils.Must(s.encrypt.DecryptString(user.Phone)),
		Socials: slices.AppendSeq(
			make([]dto.Social, 0),
			utils.Map(user.Socials, func(social models.SocialModel) dto.Social {
				return dto.Social{
					Id:    social.Id,
					Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
					Title: utils.Must(s.encrypt.DecryptString(social.Title)),
					URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
				}
			}),
		),
		Identification:             identification,
		CreatedAt:                  utils.FormatDate(user.CreatedAt),
		DateVerified:               utils.FormatDate(user.VerifiedAt.String),
		IsRestricted:               user.Restricted,
		IsBanned:                   user.Banned,
		HasSubmittedIdentification: user.HasSubmittedIdentification,
	}

	return data, nil
}

func (s *Service) GetUserById(userId string) (*response.DetailedUser, error) {
	user, err := s.userStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	isVendor, err := s.userStore.IsVendor(user.Id)
	if err != nil {
		return nil, err
	}

	if err := s.userStore.LiftRestrictionIfExpired(user.Id); err != nil {
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

		for _, expertise := range expertises {
			vendorExpertises = append(vendorExpertises, response.Expertise{
				Id:    expertise.Id,
				Title: expertise.Title,
			})
		}
	}

	response := &response.DetailedUser{
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
	}

	return response, nil
}

func (s *Service) GetUser(bearerToken string) (*response.DetailedUser, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	user, err := s.userStore.FindById(userId)
	if err != nil {
		return nil, err
	}

	isVendor, err := s.userStore.IsVendor(user.Id)
	if err != nil {
		return nil, err
	}

	if err := s.userStore.LiftRestrictionIfExpired(user.Id); err != nil {
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

		for _, expertise := range expertises {
			vendorExpertises = append(vendorExpertises, response.Expertise{
				Id:    expertise.Id,
				Title: expertise.Title,
			})
		}
	}

	response := &response.DetailedUser{
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
	}

	return response, nil
}

func (s *Service) IsVerified(bearerToken string) (bool, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return false, err
	}

	user, err := s.userStore.FindById(userId)
	if err != nil {
		return false, err
	}

	return user.Verified, nil
}

func (s *Service) AddSocial(bearerToken string, req *request.AddSocialPayload) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	social := &models.SocialModel{
		UserId: userId,
		Site:   utils.Must(s.encrypt.EncryptString(req.Site)),
		Title:  utils.Must(s.encrypt.EncryptString(req.Title)),
		Url:    utils.Must(s.encrypt.EncryptString(req.Url)),
	}

	id, err := s.userStore.AddSocial(social)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteSocial(bearerToken, id string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if err := s.userStore.DeleteSocial(userId, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddVendorExpertise(bearerToken, expertiseId string, file *multipart.FileHeader) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if isVendor, err := s.userStore.IsVendor(userId); err != nil {
		return err
	} else {
		if !isVendor {
			return errors.New(ERR_FORBIDDEN)
		}
	}

	if hasExpertise, err := s.vendorStore.HasExpertise(userId, expertiseId); err != nil {
		return err
	} else {
		if hasExpertise {
			return errors.New(ERR_DUPLICATE_EXPERTISE)
		}
	}

	if hasApplication, err := s.applicationStore.HasPendingApplication(userId, expertiseId); err != nil {
		return err
	} else {
		if hasApplication {
			return errors.New(ERR_DUPLICATE_APPLICATION)
		}
	}

	b, err := utils.FileToBytes(file)
	if err != nil {
		return err
	}

	cipher := utils.Must(s.encrypt.EncryptFile(b))
	fileData := fs.File{
		Data:     cipher,
		Category: fs.APPLICATION_PROOF_DIR,
	}

	url, err := s.fs.SaveFile(fileData)
	if err != nil {
		return err
	}

	supportingDocumentId, err := s.supportingImageStore.Create(url)
	if err != nil {
		return err
	}
	policeClearance, err := s.vendorStore.GetPoliceClearance(userId)
	if err != nil {
		return err
	}

	application := &models.ApplicationModel{
		ApplicantId:        userId,
		ExpertiseId:        expertiseId,
		SupportingDocument: supportingDocumentId,
		PoliceClearance:    policeClearance.Id,
	}
	if _, err := s.applicationStore.Create(application); err != nil {
		return err
	}

	return nil
}

func (s *Service) SetDBL(bearerToken, value string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if value == "" {
		return errors.New(ERR_INVALID_DBL)
	}

	dbl, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if dbl <= 0 {
		return errors.New(ERR_INVALID_DBL)
	}

	if err := s.vendorStore.SetDBL(userId, dbl); err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangeAddress(bearerToken string, req *request.ChangeAddressPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	address := &models.AddressModel{
		Address:   utils.Must(s.encrypt.EncryptString(req.Address)),
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	return s.userStore.ChangeAddress(userId, address)
}

func (s *Service) UpdatePhone(bearerToken string, req *request.UpdatePhonePayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	encryptedPhone := utils.Must(s.encrypt.EncryptString(req.Phone))
	return s.userStore.UpdatePhone(userId, encryptedPhone)
}
