package bug_report_repo

import "nearbyassist/internal/models"

type BugReportRepository interface {
	Create(data *models.BugReportModel) (string, error)

	GetAll(limit, offset int) ([]*models.BugReportModel, error)
	FindById(id string) (*models.BugReportModel, error)
}
