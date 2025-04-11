package report_user_repo

import "nearbyassist/internal/models"

type ReportUserRepository interface {
	Create(data *models.UserReportModel) (string, error)
	GetAllWithStatus(status string, limit, offset int) ([]*models.UserReportModel, error)
	FindById(id string) (*models.UserReportModel, error)
	GetAllReportedIs(userId string) ([]*models.UserReportModel, error)
	GetAllReportedBy(userId string) ([]*models.UserReportModel, error)
	GetImages(reportId string) ([]string, error)
	UpdateStatus(reportId, status string) error
}
