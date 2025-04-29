package dashboard_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	"nearbyassist/internal/service/core"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/utils"
	"slices"
)

type Service struct {
	dashboardStore dashboard_repo.DashboardRepository
	encrypt        core.Encryption
}

func NewService(store dashboard_repo.DashboardRepository, encrypt core.Encryption) *Service {
	return &Service{dashboardStore: store, encrypt: encrypt}
}

func (s *Service) GetDashbaordData() (*dto.Dashboard, error) {
	recentUsers, err := s.dashboardStore.RecentUsers()
	if err != nil {
		return nil, err
	}

	data := &dto.Dashboard{
		Users: dto.DashboardUsers{
			Total:      utils.Must(s.dashboardStore.TotalUsers()),
			Reported:   utils.Must(s.dashboardStore.TotalReported()),
			Restricted: utils.Must(s.dashboardStore.TotalRestricted()),
			Vendor:     utils.Must(s.dashboardStore.TotalVendors()),
			Recent: slices.AppendSeq(
				make([]dto.User, 0),
				utils.Map(recentUsers, func(user *models.UserModel) dto.User {
					return dto.User{
						Id:        user.Id,
						Name:      utils.Must(s.encrypt.DecryptString(user.Name)),
						Email:     utils.Must(s.encrypt.DecryptString(user.Email)),
						ImageURL:  user.ImageUrl,
						CreatedAt: utils.DateMonth(user.CreatedAt),
					}
				}),
			),
		},
		Services: dto.DashboardServices{
			Total:  utils.Must(s.dashboardStore.TotalServices()),
			Active: utils.Must(s.dashboardStore.TotalActiveServices()),
		},
		WeeklyBooking: *utils.Must(s.dashboardStore.GetBookingData()),
		SearchTrend: dto.Trend{
			Searches: searchhistory.New().GetAll(),
		},
		Application: dto.Application{
			Pending: utils.Must(s.dashboardStore.TotalPendingApplications()),
		},
		Report: dto.DashboardReports{
			Pending: utils.Must(s.dashboardStore.TotalActiveReports()),
		},
	}

	return data, nil
}
