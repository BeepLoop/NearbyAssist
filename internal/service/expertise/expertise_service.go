package expertise_service

import (
	"errors"
	"mime/multipart"
	"nearbyassist/internal/models"
	expertise_repo "nearbyassist/internal/repository/expertise"
	supportingimage_repo "nearbyassist/internal/repository/supporting_image"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/utils"
	"strings"
)

type Service struct {
	userStore            user_repo.UserRepository
	vendorStore          vendor_repo.VendorRepository
	expertiseStore       expertise_repo.ExpertiseRepository
	supportingImageStore supportingimage_repo.Repository
	fs                   fs.FileStorage
	encrypt              core.Encryption
	hash                 core.Hash
	jwt                  core.Authenticator
}

func NewService(userStore user_repo.UserRepository, vendorStore vendor_repo.VendorRepository, expertiseStore expertise_repo.ExpertiseRepository, supportingImageStore supportingimage_repo.Repository, fs fs.FileStorage, encrypt core.Encryption, hash core.Hash, jwt core.Authenticator) *Service {
	return &Service{
		userStore:            userStore,
		vendorStore:          vendorStore,
		expertiseStore:       expertiseStore,
		supportingImageStore: supportingImageStore,
		fs:                   fs,
		encrypt:              encrypt,
		hash:                 hash,
		jwt:                  jwt,
	}
}

func (s *Service) CreateExpertise(data *models.ExpertiseModel, tagsInput string) (string, error) {
	// Clean up the tags input
	if tagsInput != "" {
		csv := strings.Split(tagsInput, ",")
		for _, v := range csv {
			trimmed := strings.TrimSpace(v)
			if trimmed != "" {
				data.Tags = append(data.Tags, &models.TagModel{
					Title: strings.ToLower(trimmed),
				})
			}
		}
	}

	return s.expertiseStore.Create(data)
}

func (s *Service) AddUserExpertise(bearerToken, expertiseId string, file *multipart.FileHeader) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	if isVendor, err := s.userStore.IsVendor(userId); err != nil {
		return err
	} else {
		if !isVendor {
			return errors.New("forbidden")
		}
	}

	b, err := utils.FileToBytes(file)
	if err != nil {
		return err
	}

	cipher, err := s.encrypt.EncryptFile(b)
	if err != nil {
		return err
	}

	fileData := fs.File{
		Data:     cipher,
		Category: fs.APPLICATION_PROOF_DIR,
	}
	url, err := s.fs.SaveFile(fileData)
	if err != nil {
		return err
	}

	imageId, err := s.supportingImageStore.Create(url)
	if err != nil {
		return err
	}

	s.vendorStore.AddExpertise(userId, expertiseId, imageId)

	return nil
}

func (s *Service) GetAllExpertise() ([]*models.ExpertiseModel, error) {
	return s.expertiseStore.GetAll()
}

func (s *Service) FindExpertise(query string) (*models.ExpertiseModel, error) {
	return s.expertiseStore.FindByTitle(query)
}

func (s *Service) AddTagToExpertise(expertiseId string, data *models.TagModel) (string, error) {
	return s.expertiseStore.CreateTag(expertiseId, data)
}
