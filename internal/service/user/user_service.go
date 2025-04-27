package user_service

import (
	"errors"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/user"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	userStore       repository.UserRepository
	resourceService *resource_service.Service
	encrypt         core.Encryption
	hash            core.Hash
	jwt             core.Authenticator
}

func NewService(
	userStore repository.UserRepository,
	resourceService *resource_service.Service,
	encrypt core.Encryption,
	hash core.Hash,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:       userStore,
		resourceService: resourceService,
		encrypt:         encrypt,
		hash:            hash,
		jwt:             jwt,
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
			return dto.User{
				Id:       account.Id,
				Name:     utils.Must(s.encrypt.DecryptString(account.Name)),
				Email:    utils.Must(s.encrypt.DecryptString(account.Email)),
				ImageURL: account.ImageUrl,
				Address:  utils.Must(s.encrypt.DecryptString(account.Address.Address)),
				Phone:    utils.Must(s.encrypt.DecryptString(account.Phone)),
				Socials:  account.Socials,
				Identification: dto.Identification{
					Type:          account.Identification.Type,
					IdNumber:      utils.Must(s.encrypt.DecryptString(account.Identification.ReferenceNumber)),
					FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(account.Identification.FrontImageUrl)),
					BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(account.Identification.BackImageUrl)),
				},
				CreatedAt:    utils.FormatDate(account.CreatedAt),
				DateVerified: utils.FormatDate(account.VerifiedAt.String),
				IsRestricted: account.Restricted,
				IsBanned:     account.Banned,
			}
		}),
	)

	return data, nil
}

func (s *Service) FindByEmail(email string) (*dto.User, error) {
	emailHash := utils.Must(s.hash.Generate([]byte(email)))
	user, err := s.userStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, err
	}

	data := &dto.User{
		Id:       user.Id,
		Name:     utils.Must(s.encrypt.DecryptString(user.Name)),
		Email:    utils.Must(s.encrypt.DecryptString(user.Email)),
		ImageURL: user.ImageUrl,
		Address:  utils.Must(s.encrypt.DecryptString(user.Address.Address)),
		Phone:    utils.Must(s.encrypt.DecryptString(user.Phone)),
		Socials:  user.Socials,
		Identification: dto.Identification{
			Type:          user.Identification.Type,
			IdNumber:      utils.Must(s.encrypt.DecryptString(user.Identification.ReferenceNumber)),
			FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.FrontImageUrl)),
			BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(user.Identification.BackImageUrl)),
		},
		CreatedAt:    utils.FormatDate(user.CreatedAt),
		DateVerified: utils.FormatDate(user.VerifiedAt.String),
		IsRestricted: user.Restricted,
		IsBanned:     user.Banned,
	}

	return data, nil
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

	isRestricted, isRestrictionExpired, err := s.userStore.IsRestricted(user.Id)
	if err != nil {
		return nil, err
	}

	if err := s.userStore.LiftRestrictionIfExpired(user.Id); err != nil {
		return nil, err
	}

	decryptedSocials := make([]string, 0)
	for _, social := range user.Socials {
		decrypted, err := s.encrypt.DecryptString(social)
		if err != nil {
			return nil, err
		}

		decryptedSocials = append(decryptedSocials, decrypted)
	}
	user.Socials = decryptedSocials

	vendorExpertises := make([]response.Expertise, 0)
	if isVendor {
		expertises, err := s.userStore.GetExpertise(user.Id)
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

	response := &response.DetailedUser{
		Id:           user.Id,
		Name:         utils.Must(s.encrypt.DecryptString(user.Name)),
		Email:        utils.Must(s.encrypt.DecryptString(user.Email)),
		ImageUrl:     user.ImageUrl,
		IsVerified:   user.Verified,
		IsVendor:     isVendor,
		Address:      utils.Must(s.encrypt.DecryptString(user.Address.Address)),
		Phone:        utils.Must(s.encrypt.DecryptString(user.Phone)),
		Latitude:     user.Address.Latitude,
		Longitude:    user.Address.Longitude,
		Expertises:   vendorExpertises,
		Socials:      user.Socials,
		IsRestricted: isRestricted && !isRestrictionExpired,
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

func (s *Service) AddSocial(bearerToken, url string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	encrypted, err := s.encrypt.EncryptString(url)
	if err != nil {
		return err
	}

	social := new(models.SocialModel)
	social.UserId = userId
	social.Url = encrypted

	if err := s.userStore.AddSocial(social); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteSocial(bearerToken, url string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	socials, err := s.userStore.GetSocials(userId)
	if err != nil {
		return err
	}

	var socialId string
	for _, social := range socials {
		decrypted, err := s.encrypt.DecryptString(social.Url)
		if err != nil {
			return err
		}

		if decrypted == url {
			socialId = social.Id
		}
	}

	if socialId == "" {
		return errors.New("social not found")
	}

	if err := s.userStore.DeleteSocial(userId, socialId); err != nil {
		return err
	}

	return nil
}
