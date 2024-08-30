package models

import (
	"nearbyassist/internal/id_generator"

	"github.com/jmoiron/sqlx"
)

const (
	MODEL_INIT_ERROR = "Error initializing model"
)

type Model struct {
	Id        string `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"createdAt"`

	Conn        *sqlx.DB                 `json:"-" db:"-"`
	IdGenerator id_generator.IdGenerator `json:"-" db:"-"`
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt" db:"updatedAt"`
}

type SearchParams struct {
	Latitude  float64
	Longitude float64
	Radius    float64
	Query     []string
}
