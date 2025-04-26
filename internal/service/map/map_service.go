package map_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	tag_repo "nearbyassist/internal/repository/tag"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	serviceStore service_repo.ServiceRepository
	tagStore     tag_repo.TagRepository
	encrypt      core.Encryption
}

func NewService(serviceStore service_repo.ServiceRepository, tagStore tag_repo.TagRepository, encrypt core.Encryption) *Service {
	return &Service{
		serviceStore: serviceStore,
		tagStore:     tagStore,
		encrypt:      encrypt,
	}
}

func (s *Service) GetMapData(query string) (*dto.MapData, error) {
	data := &dto.MapData{
		Services: make([]dto.Service, 0),
		Tags:     make([]string, 0),
	}

	if tags, err := s.tagStore.FindAll(); err != nil {
		return data, err
	} else {
		data.Tags = slices.AppendSeq(
			make([]string, 0),
			utils.Map(tags, func(tag *models.TagModel) string { return tag.Title }),
		)
	}

	if services, err := s.serviceStore.GetAllWithTag(query); err != nil {
		return data, err
	} else {
		data.Services = slices.AppendSeq(
			make([]dto.Service, 0),
			utils.Map(services, func(service *models.ServiceModel) dto.Service {
				return dto.Service{
					Id:          service.Id,
					VendorId:    service.VendorId,
					Title:       utils.Must(s.encrypt.DecryptString(service.Title)),
					Description: utils.Must(s.encrypt.DecryptString(service.Description)),
					Rate:        utils.StringToFloat64ElseZero(service.Rate),
					Tags: slices.AppendSeq(
						make([]string, 0),
						utils.Map(service.Tags, func(tag *models.TagModel) string {
							return tag.Title
						}),
					),
					Extras: slices.AppendSeq(
						make([]dto.Extra, 0),
						utils.Map(service.Extras, func(extra *models.ExtraModel) dto.Extra {
							return dto.Extra{
								Id:          extra.Id,
								Title:       utils.Must(s.encrypt.DecryptString(extra.Title)),
								Description: utils.Must(s.encrypt.DecryptString(extra.Description)),
								Price:       extra.Price,
							}
						}),
					),
					Images: slices.AppendSeq(
						make([]dto.Image, 0),
						utils.Map(service.Images, func(image *models.ServicePhotoModel) dto.Image {
							return dto.Image{Id: image.Id, URL: image.Url}
						}),
					),
					Address: dto.Address{
						Address:   utils.Must(s.encrypt.DecryptString(service.Address.Address)),
						Latitude:  service.Address.Latitude,
						Longitude: service.Address.Longitude,
					},
					CreatedAt: utils.FormatDate(service.CreatedAt),
					UpdatedAt: utils.FormatDate(service.UpdatedAt),
				}
			}),
		)
	}

	return data, nil
}
