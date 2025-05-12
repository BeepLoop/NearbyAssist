package models

type ExtraModel struct {
	Model
	UpdateableModel
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description" validate:"required"`
	Price       string `json:"price" db:"price" validate:"required"`
	Deleted     bool   `json:"deleted" db:"deleted"`
	ServiceId   string `db:"serviceId" json:"serviceId,omitempty"`
}
