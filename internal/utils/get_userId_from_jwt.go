package utils

import (
	"errors"
	"nearbyassist/internal/authenticator"
)

func GetUserIdFromJWT(jwtSigner authenticator.Authenticator, authHeader string) (int, error) {
	token := authHeader[len("Bearer "):]

	claims, err := jwtSigner.GetClaims(token)
	if err != nil {
		return 0, err
	}

	// Check if jwt has role, if so, user is an admin or staff
	if _, ok := claims["role"].(string); ok {
		adminId, ok := claims["adminId"].(float64)
		if !ok {
			return 0, errors.New("Value adminId not found in claims")
		}

		return int(adminId), nil
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("Value userId not found in claims")
	}

	return int(userId), nil
}

func GetUserIdFromJWTString(jwtSigner authenticator.Authenticator, token string) (int, error) {
	claims, err := jwtSigner.GetClaims(token)
	if err != nil {
		return 0, err
	}

	// Check if jwt has role, if so, user is an admin or staff
	if _, ok := claims["role"].(string); ok {
		adminId, ok := claims["adminId"].(float64)
		if !ok {
			return 0, errors.New("Value adminId not found in claims")
		}

		return int(adminId), nil
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("Value userId not found in claims")
	}

	return int(userId), nil
}
