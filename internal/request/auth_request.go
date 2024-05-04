package request

type AdminLogin struct {
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}

type UserLogin struct {
	Name  string `json:"name" db:"name" validate:"required"`
	Email string `json:"email" db:"email" validate:"required"`
	Image string `json:"image" db:"image" validate:"required"`
}

type RefreshToken struct {
	Token string `json:"token" validate:"required"`
}

type Logout struct {
	Token string `json:"token" validate:"required"`
}
