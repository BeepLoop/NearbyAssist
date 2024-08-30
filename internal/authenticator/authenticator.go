package authenticator

import (
	"nearbyassist/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN_ERR  = "Error generating access token"
	REFRESH_TOKEN_ERR = "Error generating refresh token"
)

type AdminOptions struct {
	Id       string
	Username string
	Role     models.AdminRole
}

type UserOptions struct {
	Id    string
	Name  string
	Email string
}

type Authenticator interface {
	GenerateAdminAccessToken(options AdminOptions) (string, error)
	GenerateUserAccessToken(options UserOptions) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) error
	GetClaims(tokenString string) (jwt.MapClaims, error)
}
