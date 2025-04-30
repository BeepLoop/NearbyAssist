package service_service

import (
	"database/sql"
	"errors"
	"math"
	"mime/multipart"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/route_engine"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/utils"
	"slices"
	"strings"
)

const (
	ERR_FORBIDDEN_ACTION  = "action not allowed"
	ERR_NOT_FOUND         = "resource not found"
	ERR_UNAUTHORIZED      = "unauthorized"
	ERR_DUPLICATE_SERVICE = "duplicate service"
	ERR_ACTIVELY_USED     = "resource is actively in use"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	vendorStore  vendor_repo.VendorRepository
	encrypt      core.Encryption
	hash         core.Hash
	jwt          core.Authenticator
	suggest      suggestion_engine.Engine
	route        route_engine.Engine
	fs           fs.FileStorage
}

func NewService(serviceStore service_repo.ServiceRepository, vendorStore vendor_repo.VendorRepository, encrypt core.Encryption, hash core.Hash, jwt core.Authenticator, suggest suggestion_engine.Engine, route route_engine.Engine, fs fs.FileStorage) *Service {
	return &Service{
		serviceStore: serviceStore,
		vendorStore:  vendorStore,
		encrypt:      encrypt,
		hash:         hash,
		jwt:          jwt,
		suggest:      suggest,
		route:        route,
		fs:           fs,
	}
}

func (s *Service) CreateService(req *request.AddServicePayload) (string, error) {
	isVendor, err := s.serviceStore.IsVendor(req.VendorId)
	if err != nil {
		return "", err
	}
	if !isVendor {
		return "", errors.New(ERR_UNAUTHORIZED)
	}

	// Compute signature
	signature := computeSignature(req.VendorId, req.Title, req.Description, s.hash.Generate)

	service, err := s.serviceStore.FindBySignature(signature)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return "", err
	}
	if service != nil {
		return "", errors.New(ERR_DUPLICATE_SERVICE)
	}

	newService := &models.ServiceModel{
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Rate:         req.Rate,
		TagsAsString: req.Tags,
		Extras: slices.AppendSeq(
			make([]*models.ExtraModel, 0),
			utils.Map(req.Extras, func(x request.NewExtra) *models.ExtraModel {
				return &models.ExtraModel{
					Title:       utils.Must(s.encrypt.EncryptString(x.Title)),
					Description: utils.Must(s.encrypt.EncryptString(x.Description)),
					Price:       x.Price,
				}
			}),
		),
		Signature: signature,
	}

	return s.serviceStore.Create(newService)
}

func (s *Service) GetService(serviceId string) (*response.DetailedServiceResponse, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	reviews, err := s.serviceStore.GetReviews(serviceId)
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

	response := &response.DetailedServiceResponse{
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
	}

	return response, nil
}

func (s *Service) UpdateService(bearerToken string, req *request.UpdateServicePayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	// Validate that vendor of service is the current user
	vendor, err := s.vendorStore.FindById(userId)
	if err != nil {
		return err
	}
	if vendor.VendorId != req.VendorId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	updatedService := &models.ServiceModel{
		Model:        models.Model{Id: req.Id},
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Rate:         req.Rate,
		TagsAsString: req.Tags,
		Signature:    computeSignature(req.VendorId, req.Title, req.Description, s.hash.Generate),
	}

	if err := s.serviceStore.Update(updatedService); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddImage(bearerToken, serviceId string, files []*multipart.FileHeader) (*models.ServicePhotoModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	if service.VendorId != userId {
		return nil, errors.New(ERR_UNAUTHORIZED)
	}

	file := files[0]
	bytes, err := utils.FileToBytes(file)
	if err != nil {
		return nil, err
	}

	fileData := fs.File{
		Data:     bytes,
		Category: fs.SERVICE_PHOTO_DIR,
	}
	url, err := s.fs.SaveFile(fileData)
	if err != nil {
		return nil, err
	}

	photoData := &models.ServicePhotoModel{
		ServiceId: serviceId,
		VendorId:  userId,
		Url:       url,
	}

	imageId, err := s.serviceStore.AddImage(photoData)
	if err != nil {
		return nil, err
	}

	photoData.Id = imageId

	return photoData, nil
}

func (s *Service) DeleteImage(bearerToken, imageId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	image, err := s.serviceStore.FindPhotoById(imageId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(ERR_NOT_FOUND)
		}

		return err
	}
	if image.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if err := s.fs.DeleteFile(image.Url); err != nil {
		return err
	}

	if err := s.serviceStore.DeleteImage(imageId); err != nil {
		return err
	}

	return nil
}

func (s *Service) Disable(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return err
	}
	if service.VendorId != userId {
		return errors.New(ERR_FORBIDDEN_ACTION)
	}

	return s.serviceStore.Disable(service.Id)
}

func (s *Service) Enable(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return err
	}
	if service.VendorId != userId {
		return errors.New(ERR_FORBIDDEN_ACTION)
	}

	return s.serviceStore.Enable(service.Id)
}

func (s *Service) AddExtra(bearerToken string, input *request.AddExtraPayload) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	service, err := s.serviceStore.FindById(input.ServiceId)
	if err != nil {
		return "", err
	}

	if service.VendorId != userId {
		return "", errors.New(ERR_UNAUTHORIZED)
	}

	data := &models.ExtraModel{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		ServiceId:   input.ServiceId,
	}

	if encrypted, err := s.encrypt.EncryptString(data.Title); err != nil {
		return "", err
	} else {
		data.Title = encrypted
	}

	if encrypted, err := s.encrypt.EncryptString(data.Description); err != nil {
		return "", err
	} else {
		data.Description = encrypted
	}

	extraId, err := s.serviceStore.AddExtra(data)
	if err != nil {
		return "", err
	}

	return extraId, nil
}

func (s *Service) EditExtra(bearerToken string, data *request.EditExtraPayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	extra, err := s.serviceStore.FindExtraById(data.Id)
	if err != nil {
		return err
	}

	service, err := s.serviceStore.FindById(extra.ServiceId)
	if err != nil {
		return err
	}

	if service.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	updatedExtra := &models.ExtraModel{
		Model:       models.Model{Id: data.Id},
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
	}

	if encrypted, err := s.encrypt.EncryptString(updatedExtra.Title); err != nil {
		return err
	} else {
		updatedExtra.Title = encrypted
	}

	if encrypted, err := s.encrypt.EncryptString(updatedExtra.Description); err != nil {
		return err
	} else {
		updatedExtra.Description = encrypted
	}

	if err := s.serviceStore.EditExtra(updatedExtra); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteExtra(bearerToken, extraId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	extra, err := s.serviceStore.FindExtraById(extraId)
	if err != nil {
		return err
	}

	service, err := s.serviceStore.FindById(extra.ServiceId)
	if err != nil {
		return err
	}

	if service.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	isActive, err := s.serviceStore.HasActiveBookingWithThisExtra(extraId)
	if err != nil {
		return err
	}
	if isActive {
		return errors.New(ERR_ACTIVELY_USED)
	}

	return s.serviceStore.DeleteExtra(extraId)
}

func (s *Service) SearchService(params map[string]string) ([]*response.ServiceSearchResult, error) {
	// Update search history
	if q, ok := params["q"]; ok {
		tags := strings.Split(q, ",")

		for _, tag := range tags {
			cleaned := strings.ReplaceAll(tag, "_", " ")
			searchhistory.New().Insert(cleaned)
		}
	}

	origin := new(models.GeoSpatialModel)
	if location, ok := params["l"]; ok {
		if err := origin.FromString(location); err != nil {
			return nil, err
		}
	}

	tags := make([]string, 0)
	if q, ok := params["q"]; ok {
		cleaned := utils.Map(strings.Split(q, ","), func(tag string) string {
			return strings.ReplaceAll(tag, "_", " ")
		})
		slices.AppendSeq(tags, cleaned)
	}

	matchedServices, err := s.serviceStore.FuzzyMatchTags(tags)
	if err != nil {
		return nil, err
	}

	// Filter out services where vendor is banned or restricted
	validServices := slices.AppendSeq(
		make([]*models.ServiceModel, 0),
		utils.Retain(matchedServices, func(service *models.ServiceModel) bool {
			restricted, err := s.serviceStore.IsVendorRestricted(service.Id)
			if err != nil || restricted {
				return false
			}

			banned, err := s.serviceStore.IsVendorBanned(service.Id)
			if err != nil || banned {
				return false
			}

			return !restricted && !banned
		}),
	)

	// Retrieve vendor details of each service
	for _, service := range validServices {
		if vendor, err := s.vendorStore.FindById(service.VendorId); err != nil {
			return nil, err
		} else {
			service.Vendor = *vendor
		}
	}

	// Compute distance from origin to each service
	suggestionOpsInput := make([]dto.GeospatialOperation, 0)
	for _, service := range validServices {
		destination := &models.GeoSpatialModel{
			Latitude:  service.Address.Latitude,
			Longitude: service.Address.Longitude,
		}

		distance, err := s.route.GetDistance(origin, destination)
		if err != nil {
			distance = math.MaxFloat32
		}

		completedBookings, err := s.vendorStore.CompletedBookingCountOfService(service.VendorId, service.Id)
		if err != nil {
			// If error, skip this service
			continue
		}

		suggestionOpsInput = append(suggestionOpsInput, dto.GeospatialOperation{
			Id:                 service.Id,
			Rate:               utils.StringToFloat32ElseZero(service.Rate),
			Rating:             utils.StringToFloat32ElseZero(service.Vendor.Rating),
			Latitude:           float32(service.Vendor.User.Address.Latitude),
			Longitude:          float32(service.Vendor.User.Address.Longitude),
			CompletedBookings:  float32(completedBookings),
			DistanceFromOrigin: distance,
		})
	}

	// Compute suggestibility score of each service
	serviceScores, err := s.suggest.GenerateSuggestions(suggestionOpsInput)
	if err != nil {
		return nil, err
	}

	results := slices.AppendSeq(
		make([]*response.ServiceSearchResult, 0),
		utils.Map(validServices, func(service *models.ServiceModel) *response.ServiceSearchResult {
			// Index should be guaranteed to not be -1, something is terribly wrong if it is -1
			index := slices.IndexFunc(suggestionOpsInput, func(s dto.GeospatialOperation) bool {
				return s.Id == service.Id
			})

			return &response.ServiceSearchResult{
				Id:                service.Id,
				VendorName:        utils.Must(s.encrypt.DecryptString(service.Vendor.User.Name)),
				SuggestionScore:   float32(serviceScores[service.Id]),
				Rate:              float32(utils.StringToFloat64ElseZero(service.Rate)),
				Rating:            float32(utils.StringToFloat64ElseZero(service.Vendor.Rating)),
				Latitude:          service.Vendor.User.Address.Latitude,
				Longitude:         service.Vendor.User.Address.Longitude,
				CompletedBookings: suggestionOpsInput[index].CompletedBookings,
			}
		}),
	)

	return results, nil
}

func (s *Service) FindRoute(serviceId string, origin string) (route_engine.PolylineCode, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return "", err
	}

	originLatitude, originLongitude, err := models.ParseCoordinate(origin)
	if err != nil {
		return "", err
	}

	from := &models.GeoSpatialModel{
		Latitude:  originLatitude,
		Longitude: originLongitude,
	}
	distination := &models.GeoSpatialModel{
		Latitude:  service.Address.Latitude,
		Longitude: service.Address.Longitude,
	}

	polyline, err := s.route.GetPolyline(from, distination)
	if err != nil {
		return "", err
	}

	return polyline, nil
}
