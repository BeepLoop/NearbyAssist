package report_user_repo

import "nearbyassist/internal/models"

type ReportUserRepository interface {
	Create(data *models.ReportedUserModel) (string, error)
	GetAll(limit, offset int) ([]*models.ReportedUserModel, error)
	FindById(id string) (*models.ReportedUserModel, error)
	FindByUserId(id string) (*models.ReportedUserModel, error)
	CloseReport(id string) error
}
