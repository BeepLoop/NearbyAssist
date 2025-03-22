package policeclearance_repo

import "nearbyassist/internal/models"

type Repository interface {
	Create(url string) (string, error)
	FindById(id string) (*models.PoliceClearanceModel, error)
}
