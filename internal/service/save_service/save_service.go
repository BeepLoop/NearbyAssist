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

		reviews, err := s.serviceStore.GetReviews(saved.ServiceId)
		if err != nil {
			return nil, err
		}

		ratings := make([]int, 5)
		for _, review := range reviews {
			switch review.Rating {
			case 5:
				ratings[4]++
			case 4:
				ratings[3]++
			case 3:
				ratings[2]++
			case 2:
				ratings[1]++
			case 1:
				ratings[0]++
			}
		}

		vendor, err := s.vendorStore.FindById(service.VendorId)
		if err != nil {
			return nil, err
		}

		savedServices = append(savedServices, &response.DetailedServiceResponse{
			Service: response.Service{
				Id:          service.Id,
				VendorId:    service.VendorId,
				Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
				Description: utils.Must(s.encrypt.DecryptString(service.Description)),
				Price:       service.Price,
				PricingType: string(service.PricingType),
				Tags: slices.AppendSeq(
					make([]string, 0),
					utils.Map(service.Tags, func(t *models.TagModel) string {
						return t.Title
					}),
				),
				Extras: slices.AppendSeq(
					make([]response.Extra, 0),
					utils.Map(service.Extras, func(x *models.ExtraModel) response.Extra {
						return response.Extra{
							Id:          x.Id,
							Title:       utils.Must(s.encrypt.DecryptString(x.Title)),
							Description: utils.Must(s.encrypt.DecryptString(x.Description)),
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
					Latitude:  service.Address.Latitude,
					Longitude: service.Address.Longitude,
				},
				Disabled: service.Disabled,
			},
			Vendor: response.Vendor{
				Id:       vendor.VendorId,
				Name:     utils.Must(s.encrypt.DecryptString(vendor.User.Name)),
				Email:    utils.Must(s.encrypt.DecryptString(vendor.User.Email)),
				ImageUrl: vendor.User.ImageUrl,
				Phone:    vendor.User.Phone,
				Rating:   vendor.Rating,
				Socials: slices.AppendSeq(
					make([]response.Social, 0),
					utils.Map(vendor.User.Socials, func(social models.SocialModel) response.Social {
						return response.Social{
							Id:    social.Id,
							Site:  utils.Must(s.encrypt.DecryptString(social.Site)),
							Title: utils.Must(s.encrypt.DecryptString(social.Title)),
							URL:   utils.Must(s.encrypt.DecryptString(social.Url)),
						}
					}),
				),
				Expertise: slices.AppendSeq(
					make([]string, 0),
					utils.Map(vendor.Expertise, func(e models.ExpertiseModel) string { return e.Title }),
				),
				Address: vendor.User.Address.Address,
			},
			Reviews: slices.AppendSeq(
				make([]response.Review, 0),
				utils.Map(reviews, func(review *models.ReviewModel) response.Review {
					return response.Review{
						Id:               review.Id,
						BookingId:        review.BookingId,
						Rating:           review.Rating,
						Text:             utils.Must(s.encrypt.DecryptString(review.Text)),
						CreatedAt:        review.CreatedAt,
						RevieweeName:     utils.Must(s.encrypt.DecryptString(review.Reviewee.Name)),
						RevieweeImageUrl: review.Reviewee.ImageUrl,
					}
				}),
			),
			Ratings: ratings,
		})
	}

	return savedServices, nil
}
