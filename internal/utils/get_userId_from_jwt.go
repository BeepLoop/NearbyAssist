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

	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("Value userId not found in claims")
	}

	return int(userId), nil
}
