package authenticator

import (
	"nearbyassist/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateAdminAccessToken(admin *models.AdminModel) (string, error)
	GenerateUserAccessToken(user *models.UserModel) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) error
	GetClaims(tokenString string) (jwt.MapClaims, error)
}
