package auth

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

func GenerateAdminAccessToken(options AdminOptions) (string, error) {
	return "", nil
}

func GenerateUserAccessToken(options UserOptions) (string, error) {
	return "", nil
}

func GenerateRefreshToken() (string, error) {
	return "", nil
}

func ValidateToken(tokenString string) error {
	return nil
}

func GetClaims(tokenString string) (jwt.MapClaims, error) {
	return nil, nil
}
