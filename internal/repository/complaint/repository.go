package complaint_repo

import "nearbyassist/internal/models"

type ComplaintRepository interface {
	CreateSystemComplaint(data *models.SystemComplaintModel) (string, error)

	GetAll(limit, offset int) ([]*models.ComplaintModel, error)
	FindById(id string) (*models.ComplaintModel, error)
}
