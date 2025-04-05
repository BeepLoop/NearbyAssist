package vendor_service

import (
	"database/sql"
	"nearbyassist/internal/models"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	vendorStore  vendor_repo.VendorRepository
	serviceStore service_repo.ServiceRepository
	encrypt      core.Encryption
	hash         core.Hash
}

func NewService(vendorStore vendor_repo.VendorRepository, serviceStore service_repo.ServiceRepository, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		vendorStore:  vendorStore,
		serviceStore: serviceStore,
		encrypt:      encrypt,
		hash:         hash,
	}
}

func (s *Service) GetAll(limit, offset int) ([]*models.VendorModel, error) {
	accounts, err := s.vendorStore.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if decrypted, err := s.encrypt.DecryptString(account.Name); err != nil {
			return nil, err
		} else {
			account.Name = decrypted
		}

		if decrypted, err := s.encrypt.DecryptString(account.Email); err != nil {
			return nil, err
		} else {
			account.Email = decrypted
		}
	}

	return accounts, nil
}

func (s *Service) FindByEmail(email string) (*models.VendorModel, error) {
	emailHash, err := s.hash.Generate([]byte(email))
	if err != nil {
		return nil, err
	}

	vendor, err := s.vendorStore.FindByEmailHash(emailHash)
	if err != nil {
		return nil, err
	}

	if plainText, err := s.encrypt.DecryptString(vendor.Name); err != nil {
		return nil, err
	} else {
		vendor.Name = plainText
	}

	if plainText, err := s.encrypt.DecryptString(vendor.Email); err != nil {
		return nil, err
	} else {
		vendor.Email = plainText
	}

	if vendor.Phone.Valid {
		if plainText, err := s.encrypt.DecryptString(vendor.Phone.String); err != nil {
			return nil, err
		} else {
			vendor.Phone = sql.NullString{String: plainText, Valid: true}
		}
	}

	decryptedSocials := make([]string, 0)
	for _, social := range vendor.Socials {
		decrypted, err := s.encrypt.DecryptString(social)
		if err != nil {
			return nil, err
		}

		decryptedSocials = append(decryptedSocials, decrypted)
	}
	vendor.Socials = decryptedSocials

	return vendor, nil
}

func (s *Service) FindById(id string) (*models.VendorModel, error) {
	vendor, err := s.vendorStore.FindById(id)
	if err != nil {
		return nil, err
	}

	vendor.Name = utils.Must(s.encrypt.DecryptString(vendor.Name))
	vendor.Email = utils.Must(s.encrypt.DecryptString(vendor.Email))
	if vendor.Phone.Valid {
		vendor.Phone.String = utils.Must(s.encrypt.DecryptString(vendor.Phone.String))
	}

	decryptedSocials := make([]string, 0)
	for _, social := range vendor.Socials {
		decrypted, err := s.encrypt.DecryptString(social)
		if err != nil {
			return nil, err
		}

		decryptedSocials = append(decryptedSocials, decrypted)
	}
	vendor.Socials = decryptedSocials

	return vendor, nil
}

func (s *Service) GetVendorServicesList(vendorId string) ([]*models.ServiceModel, error) {
	services, err := s.vendorStore.GetVendorServiceList(vendorId)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if images, err := s.serviceStore.GetPhotos(service.Id); err != nil {
			return nil, err
		} else {
			service.Images = images
		}

		service.Title = utils.Must(s.encrypt.DecryptString(service.Title))
		service.Description = utils.Must(s.encrypt.DecryptString(service.Description))

		for _, extra := range service.Extras {
			extra.Title = utils.Must(s.encrypt.DecryptString(extra.Title))
			extra.Description = utils.Must(s.encrypt.DecryptString(extra.Description))
		}
	}

	return services, nil
}
