package vendor_service

import (
	"database/sql"
	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
)

type Service struct {
	vendorStore repository.VendorRepository
	encrypt     core.Encryption
	hash        core.Hash
}

func NewService(vendorStore repository.VendorRepository, encrypt core.Encryption, hash core.Hash) *Service {
	return &Service{
		vendorStore: vendorStore,
		encrypt:     encrypt,
		hash:        hash,
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

		if account.Phone.Valid {
			account.PhoneString = account.Phone.String
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

func (s *Service) GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error) {
	services, err := s.vendorStore.GetVendorServiceList(vendorId)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if plain, err := s.encrypt.DecryptString(service.Title); err != nil {
			return nil, err
		} else {
			service.Title = plain
		}

		if plain, err := s.encrypt.DecryptString(service.Description); err != nil {
			return nil, err
		} else {
			service.Description = plain
		}

		if tags, err := s.vendorStore.GetTags(service.Id); err != nil {
			return nil, err
		} else {
			service.Tags = tags
		}
	}

	return services, nil
}
