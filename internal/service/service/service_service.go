package service_service

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/utils"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	vendorStore  vendor_repo.VendorRepository
	encrypt      auth.Encryption
	hash         auth.Hash
	jwt          auth.Authenticator
	suggest      suggestion_engine.Engine
	route        route_engine.Engine
	fs           fs.FileStorage
}

func NewService(serviceStore service_repo.ServiceRepository, vendorStore vendor_repo.VendorRepository, encrypt auth.Encryption, hash auth.Hash, jwt auth.Authenticator, suggest suggestion_engine.Engine, route route_engine.Engine, fs fs.FileStorage) *Service {
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

func (s *Service) CreateService(req *request.NewServicePayload) (string, error) {
	if err := s.serviceStore.IsVendor(req.VendorId); err != nil {
		return "", err
	}

	encryptedTitle, err := s.encrypt.EncryptString(req.Title)
	if err != nil {
		return "", err
	}

	encryptedDesc, err := s.encrypt.EncryptString(req.Description)
	if err != nil {
		return "", err
	}

	// Compute signature
	toSign := fmt.Sprintf("%s_%s_%f_%f", req.VendorId, req.Description, req.Latitude, req.Longitude)
	signature, err := s.hash.Generate([]byte(toSign))
	if err != nil {
		return "", err
	}

	if service, err := s.serviceStore.FindBySignature(signature); err == nil && service != nil {
		return "", err
	}

	extras := make([]*models.ExtraModel, 0)
	for _, extra := range req.Extras {
		titleCipher, err := s.encrypt.EncryptString(extra.Title)
		if err != nil {
			return "", err
		}

		descCipher, err := s.encrypt.EncryptString(extra.Description)
		if err != nil {
			return "", err
		}

		extras = append(extras, &models.ExtraModel{
			Title:       titleCipher,
			Description: descCipher,
			Price:       extra.Price,
		})
	}

	newService := new(models.ServiceModel)
	newService.VendorId = req.VendorId
	newService.Title = encryptedTitle
	newService.Description = encryptedDesc
	newService.Rate = req.Rate
	newService.TagsAsString = req.Tags
	newService.Latitude = req.Latitude
	newService.Longitude = req.Longitude
	newService.Signature = signature
	newService.Extras = extras

	serviceId, err := s.serviceStore.Create(newService)
	if err != nil {
		return "", err
	}

	return serviceId, nil
}

func (s *Service) NewGetService(serviceId string) (*response.DetailedServiceResponse, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	if cipher, err := s.encrypt.DecryptString(service.Title); err != nil {
		return nil, err
	} else {
		service.Title = cipher
	}

	if cipher, err := s.encrypt.DecryptString(service.Description); err != nil {
		return nil, err
	} else {
		service.Description = cipher
	}

	for _, extra := range service.Extras {
		if cipher, err := s.encrypt.DecryptString(extra.Title); err != nil {
			return nil, err
		} else {
			extra.Title = cipher
		}

		if cipher, err := s.encrypt.DecryptString(extra.Description); err != nil {
			return nil, err
		} else {
			extra.Description = cipher
		}
	}

	reviews, err := s.serviceStore.GetReviews(serviceId)
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

	if decrypted, err := s.encrypt.DecryptString(vendor.Name); err != nil {
		return nil, err
	} else {
		vendor.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(vendor.Email); err != nil {
		return nil, err
	} else {
		vendor.Email = decrypted
	}

	if vendor.Phone.Valid {
		if decrypted, err := s.encrypt.DecryptString(vendor.Phone.String); err != nil {
			return nil, err
		} else {
			vendor.PhoneString = decrypted
		}
	}

	response := &response.DetailedServiceResponse{
		Service:        service,
		Vendor:         vendor,
		CountPerRating: countPerRating,
	}

	return response, nil
}

// NOTE: Deprecated
func (s *Service) GetService(serviceId string) (map[string]interface{}, error) {
	service, err := s.serviceStore.FindById(serviceId)
	if err != nil {
		return nil, err
	}

	if cipher, err := s.encrypt.DecryptString(service.Title); err != nil {
		return nil, err
	} else {
		service.Title = cipher
	}

	if cipher, err := s.encrypt.DecryptString(service.Description); err != nil {
		return nil, err
	} else {
		service.Description = cipher
	}

	extras := make([]*models.ExtraModel, 0)
	for _, extra := range service.Extras {
		decryptedTitle, err := s.encrypt.DecryptString(extra.Title)
		if err != nil {
			return nil, err
		}

		decryptedDescription, err := s.encrypt.DecryptString(extra.Description)
		if err != nil {
			return nil, err
		}

		extras = append(extras, &models.ExtraModel{
			Model:           extra.Model,
			UpdateableModel: extra.UpdateableModel,
			Title:           decryptedTitle,
			Description:     decryptedDescription,
			Price:           extra.Price,
		})
	}
	service.Extras = extras

	if tags, err := s.serviceStore.GetTags(serviceId); err != nil {
		return nil, err
	} else {
		service.Tags = tags
	}

	reviews, err := s.serviceStore.GetReviews(serviceId)
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

	photos, err := s.serviceStore.GetPhotos(serviceId)
	if err != nil {
		return nil, err
	}

	vendor, err := s.serviceStore.GetVendorInfo(service.VendorId)
	if err != nil {
		return nil, err
	}

	if decrypted, err := s.encrypt.DecryptString(vendor.Name); err != nil {
		return nil, err
	} else {
		vendor.Name = decrypted
	}

	if decrypted, err := s.encrypt.DecryptString(vendor.Email); err != nil {
		return nil, err
	} else {
		vendor.Email = decrypted
	}

	if vendor.Phone.Valid {
		if decrypted, err := s.encrypt.DecryptString(vendor.Phone.String); err != nil {
			return nil, err
		} else {
			vendor.Phone = sql.NullString{String: decrypted, Valid: true}
		}
	}

	vendorData := struct {
		Id           string   `json:"id"`
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		Phone        string   `json:"phone"`
		ImageUrl     string   `json:"imageUrl"`
		Rating       string   `json:"rating"`
		IsRestricted bool     `json:"isRestricted"`
		Expertise    []string `json:"expertise"`
	}{
		Id:           vendor.VendorId,
		Name:         vendor.Name,
		Email:        vendor.Email,
		Phone:        vendor.Phone.String,
		ImageUrl:     vendor.ImageUrl,
		Rating:       vendor.Rating,
		IsRestricted: vendor.Restricted,
		Expertise:    vendor.Expertise,
	}

	data := map[string]interface{}{
		"serviceInfo":    service,
		"vendorInfo":     vendorData,
		"serviceImages":  photos,
		"countPerRating": countPerRating,
	}

	return data, nil
}

func (s *Service) UpdateService(bearerToken, serviceId string, req *request.UpdateServicePayload) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if vendor, err := s.serviceStore.GetVendorInfo(userId); err != nil {
		return err
	} else {
		if vendor.VendorId != req.VendorId {
			return errors.New("unauthorized")
		}
	}

	updatedService := new(models.ServiceModel)
	updatedService.Id = req.Id
	updatedService.VendorId = req.VendorId
	updatedService.Rate = req.Rate
	updatedService.TagsAsString = req.Tags
	updatedService.Latitude = req.Latitude
	updatedService.Longitude = req.Longitude

	if cipher, err := s.encrypt.EncryptString(req.Title); err != nil {
		return err
	} else {
		updatedService.Title = cipher
	}

	if cipher, err := s.encrypt.EncryptString(req.Description); err != nil {
		return err
	} else {
		updatedService.Description = cipher
	}

	// Recompute signature
	toSign := fmt.Sprintf("%s_%s_%s_%f_%f", req.VendorId, req.Title, req.Description, req.Latitude, req.Longitude)
	if signature, err := s.hash.Generate([]byte(toSign)); err != nil {
		return err
	} else {
		updatedService.Signature = signature
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

	// NOTE: Not encrypted because its gonna be public anyway
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

	photo, err := s.serviceStore.FindPhotoById(imageId)
	if err != nil {
		return err
	}

	if photo.VendorId != userId {
		return errors.New("unauthorized")
	}

	// TODO: delete file in storage

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
