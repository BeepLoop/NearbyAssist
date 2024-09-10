package request

type AdminLoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginPayload struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Image string `json:"image" validate:"required"`
}

type TokenRefreshPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
