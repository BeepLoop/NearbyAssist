package models

type Model struct {
	Id        string `json:"id,omitempty" db:"id"`
	CreatedAt string `json:"createdAt,omitempty" db:"createdAt"`
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt,omitempty" db:"updatedAt"`
}
