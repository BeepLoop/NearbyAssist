package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIdFromToken(token string, getClaim func(string) (jwt.MapClaims, error)) (string, error) {
	claims, err := getClaim(token)
	if err != nil {
		return "", err
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("User ID not found in JWT")
	}

	return id, nil
}

func GetAdminIdFromToken(token string, getClaim func(string) (jwt.MapClaims, error)) (string, error) {
	claims, err := getClaim(token)
	if err != nil {
		return "", err
	}

	id, ok := claims["adminId"].(string)
	if !ok {
		return "", errors.New("User ID not found in JWT")
	}

	return id, nil
}
