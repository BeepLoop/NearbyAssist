package complaint

import "nearbyassist/internal/models"

type ComplaintStore interface {
	CreateSystemComplaint(data *models.SystemComplaintModel) (string, error)
}
