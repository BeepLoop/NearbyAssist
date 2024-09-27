package models

type Model struct {
	Id        string `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt" db:"updatedAt"`
}
