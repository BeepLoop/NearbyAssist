package models

type ModelOperation interface {
	Create() (int, error)
	Update(id int) error
	Delete(id int) error
}

type Locatable interface {
	GetGeolocation() (*GeoSpatialModel, error)
}

type Model struct {
	Id        int    `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt" db:"updatedAt"`
}
