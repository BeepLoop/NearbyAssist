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

		if vendor.Phone.Valid {
			vendor.Phone.String = utils.Must(s.encrypt.DecryptString(vendor.Phone.String))
		}

		savedServices = append(savedServices, &response.DetailedServiceResponse{
			Service: response.Service{
				Id:          service.Id,
				VendorId:    service.VendorId,
				Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
				Description: utils.Must(s.encrypt.DecryptString(service.Description)),
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
					Latitude:  service.Latitude,
					Longitude: service.Longitude,
				},
			},
			Vendor: response.Vendor{
				Id:        vendor.VendorId,
				Name:      utils.Must(s.encrypt.DecryptString(vendor.Name)),
				Email:     utils.Must(s.encrypt.DecryptString(vendor.Email)),
				ImageUrl:  vendor.ImageUrl,
				Phone:     vendor.Phone.String,
				Rating:    vendor.Rating,
				Socials:   vendor.Socials,
				Expertise: vendor.Expertise,
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
