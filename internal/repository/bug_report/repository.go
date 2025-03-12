package bug_report_repo

import "nearbyassist/internal/models"

type BugReportRepository interface {
	Create(data *models.BugReportModel) error

	GetAll(limit, offset int) ([]*models.BugReportModel, error)
	FindById(id int) (*models.BugReportModel, error)

	CompleteBug(id int) error
}
