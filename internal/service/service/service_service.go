package service_service

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"nearbyassist/internal/config/setting"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/geodistance"
	"nearbyassist/internal/service/route_engine"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/utils"
	"slices"
	"strings"

	"github.com/beeploop/simple-additive-weighting/roc"
	"github.com/beeploop/simple-additive-weighting/saw"
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

func NewService(
	serviceStore service_repo.ServiceRepository,
	vendorStore vendor_repo.VendorRepository,
	encrypt core.Encryption,
	hash core.Hash,
	jwt core.Authenticator,
	suggest suggestion_engine.Engine,
	route route_engine.Engine,
	fs fs.FileStorage,
) *Service {
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
	signature := computeSignature(req.VendorId, req.Title, req.Description, req.PricingType, s.hash.Generate)

	service, err := s.serviceStore.FindBySignature(signature)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return "", err
	}
	if service != nil {
		return "", errors.New(ERR_DUPLICATE_SERVICE)
	}

	tags := slices.AppendSeq(
		make([]string, 0),
		utils.Map(req.Tags, func(tag string) string {
			return strings.TrimSpace(tag)
		}),
	)

	newService := &models.ServiceModel{
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Price:        req.Price,
		PricingType:  models.PricingType(req.PricingType),
		TagsAsString: tags,
		Extras: slices.AppendSeq(
			make([]*models.ExtraModel, 0),
			utils.Map(req.Extras, func(x request.NewExtra) *models.ExtraModel {
				return &models.ExtraModel{
					Title:       utils.Must(s.encrypt.EncryptString(x.Title)),
					Description: utils.Must(s.encrypt.EncryptString(x.Description)),
					Price:       utils.Float64ToString(x.Price),
				}
			}),
		),
		Signature: signature,
	}

	var serviceId string
	if newService.PricingType == models.FIXED_PRICING {
		serviceId, err = s.serviceStore.Create(newService)
		if err != nil {
			return "", err
		}
	} else {
		serviceId, err = s.serviceStore.CreateWithPricingType(newService)
		if err != nil {
			return "", err
		}
	}

	sse.New().IncreasePendingService()

	return serviceId, nil
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
			Disabled:     service.Disabled,
			Status:       string(service.Status),
			RejectReason: service.RejectReason.String,
			CreatedAt:    service.CreatedAt,
			UpdatedAt:    service.UpdatedAt,
			AcceptedAt:   service.AcceptedAt.String,
			RejectedAt:   service.RejectedAt.String,
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

	if hasBooking, err := s.serviceStore.HasActiveBookingWithThisService(req.Id); err != nil {
		return err
	} else {
		if hasBooking {
			return errors.New(ERR_ACTIVELY_USED)
		}
	}

	tags := slices.AppendSeq(
		make([]string, 0),
		utils.Map(req.Tags, func(tag string) string {
			return strings.TrimSpace(tag)
		}),
	)

	updatedService := &models.ServiceModel{
		Model:        models.Model{Id: req.Id},
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Price:        req.Price,
		PricingType:  models.PricingType(req.PricingType),
		TagsAsString: tags,
		Signature:    computeSignature(req.VendorId, req.Title, req.Description, req.PricingType, s.hash.Generate),
	}

	if err := s.serviceStore.Update(updatedService); err != nil {
		return err
	}

	return nil
}

func (s *Service) Resubmit(bearerToken, serviceId string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return err
	}

	if service.VendorId != userId {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if service.Status != models.SERVICE_STATUS_REJECTED {
		return errors.New(ERR_FORBIDDEN_ACTION)
	}

	if err := s.serviceStore.Resubmit(serviceId); err != nil {
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
		Price:       utils.Float64ToString(input.Price),
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
		Price:       utils.Float64ToString(data.Price),
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

	tags := make([]string, 0)
	if q, ok := params["q"]; ok {
		cleaned := utils.Map(strings.Split(q, ","), func(tag string) string {
			return strings.ReplaceAll(tag, "_", " ")
		})
		tags = slices.AppendSeq(make([]string, 0), cleaned)
	}

	matchedServices := make([]*models.ServiceModel, 0)
	searchBehavior := setting.New().Values.SearchBehavior
	if searchBehavior == setting.EXACT_MATCH {
		if res, err := s.serviceStore.GetAllWithTagAny(tags); err != nil {
			return nil, err
		} else {
			matchedServices = res
		}
	} else if searchBehavior == setting.FUZZY_MATCH {
		if res, err := s.serviceStore.FuzzyMatchTags(tags); err != nil {
			return nil, err
		} else {
			matchedServices = res
		}
	} else {
		if res, err := s.serviceStore.GetAllWithTagAny(tags); err != nil {
			return nil, err
		} else {
			matchedServices = res
		}
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

	origin := new(models.GeoSpatialModel)
	if location, ok := params["l"]; ok {
		if err := origin.FromString(location); err != nil {
			return nil, err
		}
	}

	radius := 0.0
	if r, ok := params["r"]; !ok {
		return nil, errors.New("missing parameter radius")
	} else {
		radius = utils.StringToFloat64ElseZero(r)
	}

	// Filter out services outside of given radius
	inRangeServices := slices.AppendSeq(
		make([]*models.ServiceModel, 0),
		utils.Retain(validServices, func(service *models.ServiceModel) bool {
			userLocation := geodistance.Coordinate{
				Latitude:  origin.Latitude,
				Longitude: origin.Longitude,
			}

			serviceLocation := geodistance.Coordinate{
				Latitude:  service.Address.Latitude,
				Longitude: service.Address.Longitude,
			}

			distanceInMeter := userLocation.DistanceTo(serviceLocation, geodistance.M)
			if distanceInMeter > geodistance.Distance(radius) {
				return false
			}

			return true
		}),
	)

	// Retrieve vendor details of each service
	for _, service := range inRangeServices {
		if vendor, err := s.vendorStore.FindById(service.VendorId); err != nil {
			return nil, err
		} else {
			service.Vendor = *vendor
		}
	}

	// Perform SAW with ROC
	// Compute distance from origin to each service
	completedBookingMap := make(map[string]int)
	computedDistances := make(map[string]float32)
	alternatives := make([]saw.Alternative, 0)
	for _, service := range inRangeServices {
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
			fmt.Println("err completed bookins: ", err.Error())
			// If error, skip this service
			completedBookings = 0
		}

		alternative := saw.Alternative{
			Title: service.Id,
			Criterion: []saw.Criteria{
				{Title: "price", Value: utils.StringToFloat64ElseZero(service.Price)},
				{Title: "rating", Value: utils.StringToFloat64ElseZero(service.Vendor.Rating)},
				{Title: "distance", Value: float64(distance)},
				{Title: "completed_bookings", Value: float64(completedBookings)},
			},
		}

		completedBookingMap[service.Id] = completedBookings
		computedDistances[service.Id] = distance
		alternatives = append(alternatives, alternative)
	}

	userCriteriaPreference := []roc.Criteria{}
	if preferences, ok := params["preference"]; !ok {
		// Default criteria
		userCriteriaPreference = []roc.Criteria{
			{Title: "price", Rank: 1},
			{Title: "rating", Rank: 2},
			{Title: "distance", Rank: 3},
			{Title: "completed_bookings", Rank: 4},
		}
	} else {
		for i, pref := range strings.Split(preferences, ",") {
			switch pref {
			case "p":
				userCriteriaPreference = append(userCriteriaPreference, roc.Criteria{
					Title: "price", Rank: i + 1,
				})
			case "r":
				userCriteriaPreference = append(userCriteriaPreference, roc.Criteria{
					Title: "rating", Rank: i + 1,
				})
			case "d":
				userCriteriaPreference = append(userCriteriaPreference, roc.Criteria{
					Title: "distance", Rank: i + 1,
				})
			case "b":
				userCriteriaPreference = append(userCriteriaPreference, roc.Criteria{
					Title: "completed_bookings", Rank: i + 1,
				})
			}
		}
	}

	r := roc.NewRankOrderCentroid(userCriteriaPreference)
	priceWeight := r.CalculateWeightOf("price")
	ratingWeight := r.CalculateWeightOf("rating")
	distanceWeight := r.CalculateWeightOf("distance")
	cbWeight := r.CalculateWeightOf("completed_bookings")

	scores := make(map[string]float64)
	sawInstance := saw.NewSAW(alternatives)
	normalizer := saw.NewNormalizer()

	for _, alternative := range sawInstance.Alternatives {
		price, _ := alternative.CriteriaWithTitle("price")
		otherPrices := sawInstance.CriteriasWithTitle("price")
		normalizedPrice := normalizer.NormalizeCost(price.Value, otherPrices)

		rating, _ := alternative.CriteriaWithTitle("rating")
		otherRatings := sawInstance.CriteriasWithTitle("rating")
		normalizedRating := normalizer.NormalizeBenefit(rating.Value, otherRatings)

		distance, _ := alternative.CriteriaWithTitle("distance")
		otherDistances := sawInstance.CriteriasWithTitle("distance")
		normalizedDistance := normalizer.NormalizeCost(distance.Value, otherDistances)

		cb, _ := alternative.CriteriaWithTitle("completed_bookings")
		otherCBs := sawInstance.CriteriasWithTitle("completed_bookings")
		normalizedCBs := normalizer.NormalizeBenefit(cb.Value, otherCBs)

		pairs := []saw.WeightAndNormalizedPair{
			{Weight: priceWeight, Normalized: normalizedPrice},
			{Weight: ratingWeight, Normalized: normalizedRating},
			{Weight: distanceWeight, Normalized: normalizedDistance},
			{Weight: cbWeight, Normalized: normalizedCBs},
		}

		score := sawInstance.ComputeWeightedSum(pairs)
		scores[alternative.Title] = score
	}

	results := slices.AppendSeq(
		make([]*response.ServiceSearchResult, 0),
		utils.Map(inRangeServices, func(service *models.ServiceModel) *response.ServiceSearchResult {
			return &response.ServiceSearchResult{
				Id:                service.Id,
				VendorName:        utils.Must(s.encrypt.DecryptString(service.Vendor.User.Name)),
				Suggestibility:    float32(scores[service.Id]),
				Price:             service.Price,
				Rating:            service.Vendor.Rating,
				Latitude:          service.Vendor.User.Address.Latitude,
				Longitude:         service.Vendor.User.Address.Longitude,
				CompletedBookings: completedBookingMap[service.Id],
				Distance:          computedDistances[service.Id],
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
					Disabled:     service.Disabled,
					Status:       string(service.Status),
					RejectReason: utils.Must(s.encrypt.DecryptString(service.RejectReason.String)),
					CreatedAt:    service.CreatedAt,
					UpdatedAt:    service.UpdatedAt,
					AcceptedAt:   service.AcceptedAt.String,
					RejectedAt:   service.RejectedAt.String,
				},
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
