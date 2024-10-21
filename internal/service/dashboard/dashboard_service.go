package dashboard_service

import (
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	"nearbyassist/internal/response"
)

type Service struct {
	store dashboard_repo.DashboardRepository
}

func NewService(store dashboard_repo.DashboardRepository) *Service {
	return &Service{store: store}
}

func (s *Service) GetAnalytics() (*response.Analytics, error) {
	analytics := new(response.Analytics)

	if count, err := s.store.UserCount(dashboard_repo.USER_STATUS_ALL); err != nil {
		return nil, err
	} else {
		analytics.User = count
	}

	if count, err := s.store.UserCount(dashboard_repo.USER_STATUS_VERIFIED); err != nil {
		return nil, err
	} else {
		analytics.VerifiedUser = count
	}

	if count, err := s.store.VendorCount(dashboard_repo.VENDOR_STATUS_ALL); err != nil {
		return nil, err
	} else {
		analytics.Vendor = count
	}

	if count, err := s.store.ApplicationCount(dashboard_repo.APPLICATION_STATUS_PENDING); err != nil {
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
