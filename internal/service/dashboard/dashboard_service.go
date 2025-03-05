package dashboard_service

import (
	"nearbyassist/internal/models"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
)

type Service struct {
	store dashboard_repo.DashboardRepository
}

func NewService(store dashboard_repo.DashboardRepository) *Service {
	return &Service{store: store}
}

func (s *Service) GetAnalytics() (*models.DashboardModel, error) {
	userData, err := s.store.GetUserData()
	if err != nil {
		return nil, err
	}

	bugReportData, err := s.store.GetBugReportData()
	if err != nil {
		return nil, err
	}

	vendorReportData, err := s.store.GetVendorReportData()
	if err != nil {
		return nil, err
	}

	verificationRequests, err := s.store.GetIdentityVerificationRequestsData()
	if err != nil {
		return nil, err
	}

	vendorRequests, err := s.store.GetVendorApplicationRequestsData()
	if err != nil {
		return nil, err
	}

	transactionData, err := s.store.GetTransactionData()
	if err != nil {
		return nil, err
	}

	dashboardData := &models.DashboardModel{
		UserData:        *userData,
		TransactionData: *transactionData,
		ReportData: models.ReportData{
			WeeklyBugReport:    *bugReportData,
			WeeklyVendorReport: *vendorReportData,
		},
		RequestData: models.RequestData{
			IdentityVerification: *verificationRequests,
			VendorApplication:    *vendorRequests,
		},
	}

	return dashboardData, nil
}
