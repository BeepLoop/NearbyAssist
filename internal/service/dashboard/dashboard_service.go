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

	transactionData, err := s.store.GetTransactionData()
	if err != nil {
		return nil, err
	}

	dashboardData := &models.DashboardModel{
		UserData:        *userData,
		BugReportData:   *bugReportData,
		TransactionData: *transactionData,
	}

	return dashboardData, nil
}
