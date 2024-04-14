package authenticator

import (
	"nearbyassist/internal/models"
)

type Authenticator interface {
	GenerateAccessToken(user *models.UserModel) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) error
}
