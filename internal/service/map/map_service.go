package map_service

import (
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	encrypt      core.Encryption
}

func NewService(serviceStore service_repo.ServiceRepository, encrypt core.Encryption) *Service {
	return &Service{
		serviceStore: serviceStore,
		encrypt:      encrypt,
	}
}

func (s *Service) GetServices(query string) ([]*models.ServiceModel, error) {
	services, err := s.serviceStore.GetAllWithTag(query)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		service.Title = utils.Must(s.encrypt.DecryptString(service.Title))
		service.Description = utils.Must(s.encrypt.DecryptString(service.Description))
		service.Address.Address = utils.Must(s.encrypt.DecryptString(service.Address.Address))

		service.Extras = slices.AppendSeq(
			make([]*models.ExtraModel, 0),
			utils.Map(service.Extras, func(extra *models.ExtraModel) *models.ExtraModel {
				return &models.ExtraModel{
					Model:           extra.Model,
					UpdateableModel: extra.UpdateableModel,
					Title:           utils.Must(s.encrypt.DecryptString(extra.Title)),
					Description:     utils.Must(s.encrypt.DecryptString(extra.Description)),
					Price:           extra.Price,
					Deleted:         extra.Deleted,
					ServiceId:       extra.ServiceId,
				}
			}),
		)
	}

	return services, nil
}
