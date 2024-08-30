package request

import "nearbyassist/internal/models"

type NewAdminRequest struct {
	Username string           `json:"username" db:"username" validate:"required"`
	Password string           `json:"password" db:"password" validate:"required"`
	Role     models.AdminRole `json:"role" db:"role" validate:"required"`
}
