package request

type RefreshToken struct {
	Token string `json:"token" validate:"required"`
}

type Logout struct {
	Token string `json:"token" validate:"required"`
}
