package utils

import (
	"nearbyassist/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(user types.User) (string, error) {
	claims := &types.JwtClaims{
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}
