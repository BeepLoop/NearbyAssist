package request

import "nearbyassist/internal/models"

type NewComplaint struct {
	models.Model
	Code    int    `json:"code" db:"code" validate:"required"`
	Title   string `json:"title" db:"title" validate:"required"`
	Content string `json:"content" db:"content" validate:"required"`
}

type SystemComplaint struct {
	models.Model
	Title  string `json:"title" db:"title" validate:"required"`
	Detail string `json:"detail" db:"detail" validate:"required"`
}
