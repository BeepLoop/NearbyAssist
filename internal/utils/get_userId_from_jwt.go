package utils

import (
	"errors"
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIdFromJWT(jwtSigner authenticator.Authenticator, authHeader string) (string, error) {
	token := authHeader[len("Bearer "):]

	claims, err := jwtSigner.GetClaims(token)
	if err != nil {
		return "", err
	}

	// Check if jwt has role, if so, user is an admin or staff
	if _, err := GetRoleFromClaims(claims); err == nil {
		adminId, ok := claims["adminId"].(string)
		if !ok {
			return "", errors.New("Value adminId not found in claims")
		}

		return adminId, nil
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("Value userId not found in claims")
	}

	return userId, nil
}

func GetUserIdFromJwtString(jwtSigner authenticator.Authenticator, token string) (string, error) {
	claims, err := jwtSigner.GetClaims(token)
	if err != nil {
		return "", err
	}

	// Check if jwt has role, if so, user is an admin or staff
	if _, err := GetRoleFromClaims(claims); err == nil {
		adminId, ok := claims["adminId"].(string)
		if !ok {
			return "", errors.New("Value adminId not found in claims")
		}

		return adminId, nil
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("Value userId not found in claims")
	}

	return userId, nil
}

func GetRoleFromClaims(claims jwt.MapClaims) (models.AdminRole, error) {
	if role, ok := claims["role"].(string); ok {
		return models.AdminRole(role), nil
	}

	return "", errors.New("No role field found in the JWT")
}
