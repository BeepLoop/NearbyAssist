package management

import "nearbyassist/internal/models"

type ManagementStore interface {
	CreateStaff(data *models.AdminModel) (string, error)
}
