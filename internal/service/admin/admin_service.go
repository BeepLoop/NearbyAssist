package admin_service

import (
	"errors"

	"nearbyassist/internal/models"
	repository "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
)

type service struct {
	store   repository.AdminRepository
	encrypt auth.Encryption
	hash    auth.Hash
}

func NewService(store repository.AdminRepository, encrypt auth.Encryption, hash auth.Hash) *service {
	return &service{store: store, encrypt: encrypt, hash: hash}
}

func (s *service) Login(username, password string) (*models.AdminModel, error) {
	usernameHash, err := s.hash.Generate([]byte(username))
	if err != nil {
		return nil, err
	}

	admin, err := s.store.FindByUsernameHash(usernameHash)
	if err != nil {
		return nil, err
	}

	if auth.IsPasswordMatch(admin.Password, password) == false {
		return nil, errors.New("Invalid credentials")
	}

	return admin, nil
}

func (s *service) Analytics() (*response.Analytics, error) {
	analytics := new(response.Analytics)

	if count, err := s.store.UserCount(repository.USER_STATUS_ALL); err != nil {
		return nil, err
	} else {
		analytics.User = count
	}

	if count, err := s.store.UserCount(repository.USER_STATUS_VERIFIED); err != nil {
		return nil, err
	} else {
		analytics.VerifiedUser = count
	}

	if count, err := s.store.VendorCount(repository.VENDOR_STATUS_ALL); err != nil {
		return nil, err
	} else {
		analytics.Vendor = count
	}

	if count, err := s.store.ApplicationCount(repository.APPLICATION_STATUS_PENDING); err != nil {
		return nil, err
	} else {
		analytics.PendingApplication = count
	}

	if count, err := s.store.ComplaintCount(); err != nil {
		return nil, err
	} else {
		analytics.Complaint = count
	}

	return analytics, nil
}

func (s *service) Complaints() ([]models.ComplaintModel, error) {
	complaints := make([]models.ComplaintModel, 0)

	return complaints, nil
}

func (s *service) Applications() ([]models.ApplicationModel, error) {
	applications := make([]models.ApplicationModel, 0)

	return applications, nil
}

func (s *service) IdentityVerificationRequests() ([]models.IdentityVerificationModel, error) {
	requests := make([]models.IdentityVerificationModel, 0)

	return requests, nil
}
