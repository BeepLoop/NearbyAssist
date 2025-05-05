package report_user_repo

import "nearbyassist/internal/models"

type ReportUserRepository interface {
	Create(data *models.UserReportModel) error
	GetAllWithStatus(status string, limit, offset int) ([]*models.UserReportModel, error)
	FindById(id int) (*models.UserReportModel, error)
	GetAllReportedIs(userId string) ([]*models.UserReportModel, error)
	GetAllReportedBy(userId string) ([]*models.UserReportModel, error)
	UpdateStatus(reportId int, status string) error
	Close(reportId int, action, adminId, note string) error
}
