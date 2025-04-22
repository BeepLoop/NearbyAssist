package vendor_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	vendorStore     vendor_repo.VendorRepository
	serviceStore    service_repo.ServiceRepository
	resourceService *resource_service.Service
	encrypt         core.Encryption
	hash            core.Hash
}

func NewService(
	vendorStore vendor_repo.VendorRepository,
	serviceStore service_repo.ServiceRepository,
	resourceService *resource_service.Service,
	encrypt core.Encryption,
	hash core.Hash,
) *Service {
	return &Service{
		vendorStore:     vendorStore,
		serviceStore:    serviceStore,
		resourceService: resourceService,
		encrypt:         encrypt,
		hash:            hash,
	}
}

func (s *Service) GetAll(limit, offset int) ([]dto.Vendor, error) {
	vendors, err := s.vendorStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	data := slices.AppendSeq(
		make([]dto.Vendor, 0),
		utils.Map(vendors, func(vendor *models.VendorModel) dto.Vendor {
			return dto.Vendor{
				Id:       vendor.VendorId,
				Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
				Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
				ImageURL: vendor.User.ImageUrl,
				Address:  utils.Must(s.encrypt.DecryptString(vendor.User.Address.Address)),
				Phone:    utils.Must(s.encrypt.DecryptString(vendor.User.Phone)),
				Socials: slices.AppendSeq(
					make([]string, 0),
					utils.Map(vendor.User.Socials, func(social string) string {
						return utils.Must(s.encrypt.DecryptString(social))
					}),
				),
				Identification: dto.Identification{
					Type:     vendor.User.Identification.Type,
					IdNumber: vendor.User.Identification.ReferenceNumber,
				},
				Rating:       vendor.Rating,
				Expertise:    make([]dto.Expertise, 0),
				JoinedAt:     vendor.JoinedAt,
				DateVerified: utils.FormatDate(vendor.User.VerifiedAt.String),
				IsRestricted: vendor.User.Restricted,
				IsBanned:     vendor.User.Banned,
			}
		}),
	)

	return data, nil
}

func (s *Service) FindByEmail(email string) (*dto.Vendor, error) {
	emailHash, err := s.hash.Generate([]byte(email))
	if err != nil {
		return nil, err
	}

	vendor, err := s.vendorStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, err
	}

	data := &dto.Vendor{
		Id:       vendor.VendorId,
		Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
		Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
		ImageURL: vendor.User.ImageUrl,
		Address:  utils.Must(s.encrypt.DecryptString(vendor.User.Address.Address)),
		Phone:    utils.Must(s.encrypt.DecryptString(vendor.User.Phone)),
		Socials: slices.AppendSeq(
			make([]string, 0),
			utils.Map(vendor.User.Socials, func(social string) string {
				return utils.Must(s.encrypt.DecryptString(social))
			}),
		),
		Identification: dto.Identification{
			Type:          vendor.User.Identification.Type,
			IdNumber:      vendor.User.Identification.ReferenceNumber,
			FrontImageURL: utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.FrontImageUrl)),
			BackImageURL:  utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.BackImageUrl)),
		},
		Rating:       vendor.Rating,
		Expertise:    make([]dto.Expertise, 0),
		JoinedAt:     vendor.JoinedAt,
		DateVerified: utils.FormatDate(vendor.User.VerifiedAt.String),
		IsRestricted: vendor.User.Restricted,
		IsBanned:     vendor.User.Banned,
	}

	return data, nil
}

func (s *Service) FindById(id string) (*models.VendorModel, error) {
	vendor, err := s.vendorStore.FindById(id)
	if err != nil {
		return nil, err
	}

	data := &models.VendorModel{
		VendorId: vendor.VendorId,
		Rating:   vendor.Rating,
		JoinedAt: vendor.JoinedAt,
		User: models.UserModel{
			Model:      vendor.User.Model,
			Name:       utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
			Email:      utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
			ImageUrl:   vendor.User.ImageUrl,
			Phone:      utils.Must(s.encrypt.DecryptString(vendor.User.Phone)),
			Verified:   vendor.User.Verified,
			VerifiedAt: vendor.User.VerifiedAt,
			Banned:     vendor.User.Banned,
			Restricted: vendor.User.Restricted,
			Socials: slices.AppendSeq(
				make([]string, 0),
				utils.Map(vendor.User.Socials, func(social string) string {
					return utils.Must(s.encrypt.DecryptString(social))
				}),
			),
			Address: models.AddressModel{
				Id:        vendor.User.Address.Id,
				Address:   utils.Must(s.encrypt.DecryptString(vendor.User.Address.Address)),
				Latitude:  vendor.User.Address.Latitude,
				Longitude: vendor.User.Address.Longitude,
			},
			Identification: models.IdentificationModel{
				Id:              vendor.User.Identification.Id,
				Type:            vendor.User.Identification.Type,
				ReferenceNumber: vendor.User.Identification.ReferenceNumber,
				FrontImageUrl:   utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.FrontImageUrl)),
				BackImageUrl:    utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.BackImageUrl)),
				SelfieImageUrl:  utils.Must(s.resourceService.SignURLWithDefaultDuration(vendor.User.Identification.SelfieImageUrl)),
				CreatedAt:       vendor.User.Identification.CreatedAt,
			},
		},
		Expertise: vendor.Expertise,
	}

	return data, nil
}

func (s *Service) GetVendorServicesList(vendorId string) ([]*models.ServiceModel, error) {
	services, err := s.vendorStore.GetVendorServiceList(vendorId)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if images, err := s.serviceStore.GetPhotos(service.Id); err != nil {
			return nil, err
		} else {
			service.Images = images
		}

		service.Title = utils.Must(s.encrypt.DecryptString(service.Title))
		service.Description = utils.Must(s.encrypt.DecryptString(service.Description))

		for _, extra := range service.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return services, nil
}
