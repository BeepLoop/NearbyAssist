package save_service

import (
	"nearbyassist/internal/models"
	saved_service_repo "nearbyassist/internal/repository/saved_service"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	savedServiceStore saved_service_repo.SavedServiceRepository
	serviceStore      service_repo.ServiceRepository
	vendorStore       vendor_repo.VendorRepository
	jwt               core.Authenticator
	encrypt           core.Encryption
}

func NewService(savedServiceStore saved_service_repo.SavedServiceRepository, serviceStore service_repo.ServiceRepository, vendorStore vendor_repo.VendorRepository, jwt core.Authenticator, encrypt core.Encryption) *Service {
	return &Service{
		savedServiceStore: savedServiceStore,
		serviceStore:      serviceStore,
		vendorStore:       vendorStore,
		jwt:               jwt,
		encrypt:           encrypt,
	}
}

func (s *Service) SaveService(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	data := &models.SavedServiceModel{
		UserId:    userId,
		ServiceId: serviceId,
	}

	return s.savedServiceStore.SaveService(data)
}

func (s *Service) UnsaveService(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	data := &models.SavedServiceModel{
		UserId:    userId,
		ServiceId: serviceId,
	}

	return s.savedServiceStore.UnsaveService(data)
}

func (s *Service) GetSavedServices(bearerToken string) ([]*response.DetailedServiceResponse, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	saves, err := s.savedServiceStore.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	savedServices := make([]*response.DetailedServiceResponse, 0)

	for _, saved := range saves {
		service, err := s.serviceStore.FindById(saved.ServiceId)
		if err != nil {
			return nil, err
		}

		service.Title = utils.Must(s.encrypt.DecryptString(service.Title))
		service.Description = utils.Must(s.encrypt.DecryptString(service.Description))

		for _, extra := range service.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}

		reviews, err := s.serviceStore.GetReviews(saved.ServiceId)
		if err != nil {
			return nil, err
		}

		countPerRating := response.NewCountPerRating()
		for _, review := range reviews {
			switch review.Rating {
			case 5:
				countPerRating["five"]++
			case 4:
				countPerRating["four"]++
			case 3:
				countPerRating["three"]++
			case 2:
				countPerRating["two"]++
			case 1:
				countPerRating["one"]++
			}
		}

		vendor, err := s.vendorStore.FindById(service.VendorId)
		if err != nil {
			return nil, err
		}

		vendor.Name = utils.Must(s.encrypt.DecryptString(vendor.Name))
		vendor.Email = utils.Must(s.encrypt.DecryptString(vendor.Email))
		if vendor.Phone.Valid {
			vendor.Phone.String = utils.Must(s.encrypt.DecryptString(vendor.Phone.String))
		}

		savedServices = append(savedServices, &response.DetailedServiceResponse{
			Service: response.Service{
				Id:          service.Id,
				VendorId:    service.VendorId,
				Title:       service.Title,
				Description: service.Description,
				Rate:        service.Rate,
				Tags: slices.AppendSeq(
					make([]response.Tag, 0),
					utils.Map(service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{Id: t.Id, Title: t.Title}
					}),
				),
				Extras: slices.AppendSeq(
					make([]response.Extra, 0),
					utils.Map(service.Extras, func(x *models.ExtraModel) response.Extra {
						return response.Extra{
							Id:          x.Id,
							Title:       x.Title,
							Description: x.Description,
							Price:       x.Price,
						}
					}),
				),
				Images: slices.AppendSeq(
					make([]response.Image, 0),
					utils.Map(service.Images, func(i *models.ServicePhotoModel) response.Image {
						return response.Image{Id: i.Id, Url: i.Url}
					}),
				),
				Location: response.Location{
					Latitude:  service.Latitude,
					Longitude: service.Longitude,
				},
			},
			Vendor: response.Vendor{
				Id:        vendor.VendorId,
				Name:      vendor.Name,
				Email:     vendor.Email,
				ImageUrl:  vendor.ImageUrl,
				Phone:     vendor.Phone.String,
				Rating:    vendor.Rating,
				Socials:   vendor.Socials,
				Expertise: vendor.Expertise,
			},
			CountPerRating: countPerRating,
		})
	}

	return savedServices, nil
}
