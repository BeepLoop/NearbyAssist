package dto

type Admin struct {
	Id                 string
	Username           string
	Email              string
	Role               string
	MustChangePassword bool
	Suspended          bool
	CreatedAt          string
	UpdatedAt          string
}
