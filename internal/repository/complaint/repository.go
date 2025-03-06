package complaint_repo

import "nearbyassist/internal/models"

type ComplaintRepository interface {
	CreateBugReport(data *models.BugReportModel) (string, error)

	GetAll(limit, offset int) ([]*models.ComplaintModel, error)
	FindById(id string) (*models.ComplaintModel, error)
}
