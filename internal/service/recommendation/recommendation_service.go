package recommendation_service

import (
	service_repo "nearbyassist/internal/repository/service"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/utils"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	vendorStore  vendor_repo.VendorRepository
	encrypt      core.Encryption
}

func NewService(serviceStore service_repo.ServiceRepository, vendorStore vendor_repo.VendorRepository, encrypt core.Encryption) *Service {
	return &Service{
		serviceStore: serviceStore,
		vendorStore:  vendorStore,
		encrypt:      encrypt,
	}
}

func (s *Service) GetRecommendations(limit, offset int) (*response.Recommendation, error) {
	services, err := s.serviceStore.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	recommendServices := make([]*response.ServiceRecommendation, 0)

	for _, service := range services {
		for _, extra := range service.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}

		vendor, err := s.vendorStore.FindById(service.VendorId)
		if err != nil {
			return nil, err
		}

		var thumbnail string
		if len(service.Images) > 0 {
			thumbnail = service.Images[0].Url
		}

		recommendServices = append(recommendServices, &response.ServiceRecommendation{
			Id:          service.Id,
			VendorId:    service.VendorId,
			Vendor:      utils.Must(s.encrypt.DecryptString(vendor.Name)),
			Thumbnail:   thumbnail,
			Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
			Description: utils.Must(s.encrypt.DecryptString(service.Description)),
			Rating:      vendor.Rating,
			Rate:        service.Rate,
			Tags:        service.Tags,
		})
	}

	recommendation := &response.Recommendation{
		Searches: searchhistory.Instance.GetAll(),
		Services: recommendServices,
	}

	return recommendation, nil
}
