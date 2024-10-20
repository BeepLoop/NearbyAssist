package complaint_repo

import "nearbyassist/internal/models"

type ComplaintRepository interface {
	CreateSystemComplaint(data *models.SystemComplaintModel) (string, error)
}
