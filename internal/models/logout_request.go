package models

type RefreshToken struct {
	Token string `json:"token" validate:"required"`
}

func NewRefreshTokenModel() *RefreshToken {
	return &RefreshToken{}
}
