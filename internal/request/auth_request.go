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

type UserRegisterPayload struct {
	Name      string  `json:"name" validate:"required"`
	Email     string  `json:"email" validate:"required"`
	ImageURL  string  `json:"imageUrl" validate:"required"`
	Phone     string  `json:"phone" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type TokenRefreshPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
