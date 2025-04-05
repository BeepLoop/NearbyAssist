package service_service

import (
	"errors"
	"fmt"
	"mime/multipart"
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
	if err := s.serviceStore.IsVendor(req.VendorId); err != nil {
		return "", err
	}

	// Compute signature
	rawStr := fmt.Sprintf("%s_%s_%f_%f", req.VendorId, req.Description, req.Location.Latitude, req.Location.Longitude)
	signature := utils.Must(s.hash.Generate([]byte(rawStr)))

	service, err := s.serviceStore.FindBySignature(signature)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return "", err
	}
	if service != nil {
		return "", errors.New("duplicate service")
	}

	newService := &models.ServiceModel{
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Rate:         req.Rate,
		TagsAsString: req.Tags,
		GeoSpatialModel: models.GeoSpatialModel{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
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

	serviceId, err := s.serviceStore.Create(newService)
	if err != nil {
		return "", err
	}

	return serviceId, nil
}

func (s *Service) GetService(serviceId string) (*response.DetailedServiceResponse, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	service.Title = utils.Must(s.encrypt.DecryptString(service.Title))
	service.Description = utils.Must(s.encrypt.DecryptString(service.Description))

	for _, extra := range service.Extras {
		extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
		extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
	}

	reviews, err := s.serviceStore.GetReviews(serviceId)
	if err != nil {
		return nil, err
	}

	for _, review := range reviews {
		review.Text = utils.Must(s.encrypt.DecryptString(review.Text))
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

	response := &response.DetailedServiceResponse{
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
		return errors.New("unauthorized")
	}

	// Recompute signature
	rawStr := fmt.Sprintf("%s_%s_%s_%f_%f", req.VendorId, req.Title, req.Description, req.Location.Latitude, req.Location.Longitude)
	signature := utils.Must(s.hash.Generate([]byte(rawStr)))

	updatedService := &models.ServiceModel{
		Model:        models.Model{Id: req.Id},
		VendorId:     req.VendorId,
		Title:        utils.Must(s.encrypt.EncryptString(req.Title)),
		Description:  utils.Must(s.encrypt.EncryptString(req.Description)),
		Rate:         req.Rate,
		TagsAsString: req.Tags,
		GeoSpatialModel: models.GeoSpatialModel{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		Signature: signature,
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
		return nil, errors.New("unauthorized")
	}

	file := files[0]
	bytes, err := utils.FileToBytes(file)
	if err != nil {
		return nil, err
	}

	cipher, err := s.encrypt.EncryptFile(bytes)
	if err != nil {
		return nil, err
	}

	fileData := fs.File{
		Data:     cipher,
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
		return err
	}

	if image.VendorId != userId {
		return errors.New("unauthorized")
	}

	if err := s.fs.DeleteFile(image.Url); err != nil {
		return err
	}

	if err := s.serviceStore.DeleteImage(imageId); err != nil {
		return err
	}

	return nil
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
		return "", errors.New("unauthorized")
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
		return errors.New("unauthorized")
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
		fmt.Println("find extra by id: ", err.Error())
		return err
	}

	service, err := s.serviceStore.FindById(extra.ServiceId)
	if err != nil {
		fmt.Println("find service by id: ", err.Error())
		return err
	}

	if service.VendorId != userId {
		return errors.New("unauthorized")
	}

	if err := s.serviceStore.DeleteExtra(extraId); err != nil {
		fmt.Println("delete extra: ", err.Error())
		return err
	}

	return nil
}

func (s *Service) SearchService(params map[string]string) ([]*response.ServiceSearchResult, error) {
	// Update search history
	if q, ok := params["q"]; ok {
		tags := strings.Split(q, ",")

		for _, tag := range tags {
			cleaned := strings.ReplaceAll(tag, "_", " ")
			searchhistory.Instance.Insert(cleaned)
		}
	}

	services, err := s.serviceStore.GeoSpatialSearch(params)
	if err != nil {
		return nil, err
	}

	// Filter out services with restricted OR banned vendor
	validServices := make([]*models.GeoSpatialSearchResult, 0)
	for _, service := range services {
		restricted, err := s.serviceStore.IsVendorRestricted(service.Id)
		if err != nil {
			return nil, err
		}

		banned, err := s.serviceStore.IsVendorBanned(service.Id)
		if err != nil {
			return nil, err
		}

		if !restricted && !banned {
			validServices = append(validServices, service)
		}
	}

	// Compute service distance
	for _, service := range validServices {
		origin := new(models.GeoSpatialModel)
		if location, ok := params["l"]; ok {
			if err := origin.FromString(location); err != nil {
				return nil, err
			}
		}

		destination := new(models.GeoSpatialModel)
		destination.Latitude = service.Latitude
		destination.Longitude = service.Longitude

		if distance, err := s.route.GetDistance(origin, destination); err != nil {
			return nil, err
		} else {
			service.Distance = distance
		}
	}

	// Decrypt vendor name
	for _, service := range services {
		decrypted, err := s.encrypt.DecryptString(service.VendorName)
		if err != nil {
			return nil, err
		}

		service.VendorName = decrypted
	}

	// Compute service suggestion score
	scoredServices, err := s.suggest.GenerateSuggestions(services)
	if err != nil {
		return nil, err
	}

	return scoredServices, nil
}

func (s *Service) FindRoute(serviceId string, origin string) (route_engine.PolylineCode, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return "", err
	}

	lat, long, err := models.ParseCoordinate(origin)
	if err != nil {
		return "", err
	}

	from := &models.GeoSpatialModel{
		Latitude:  lat,
		Longitude: long,
	}
	distination := &models.GeoSpatialModel{
		Latitude:  service.Latitude,
		Longitude: service.Longitude,
	}

	polyline, err := s.route.GetPolyline(from, distination)
	if err != nil {
		return "", err
	}

	return polyline, nil
}
