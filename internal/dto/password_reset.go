package dto

type PasswordResetRequest struct {
	Id        string
	AdminId   string
	Username  string
	Email     string
	CreatedAt string
}
