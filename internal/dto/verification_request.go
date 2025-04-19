package dto

type VerificationRequest struct {
	Id              string
	UserID          string
	Name            string
	Email           string
	ImageURL        string
	Phone           string
	Address         string
	IDType          string
	ReferenceNumber string
	IDFrontImageURL string
	IDBackImageURL  string
	SelfieImageURL  string
	CreatedAt       string
}
