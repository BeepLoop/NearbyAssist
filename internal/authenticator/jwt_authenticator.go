package authenticator

import (
	"errors"
	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtAuthenticator struct {
	secret        string
	signMethod    jwt.SigningMethod
	tokenDuration time.Duration
}

func NewJWTAuthenticator(conf *config.Config) *jwtAuthenticator {
	return &jwtAuthenticator{
		secret:        conf.JwtSecret,
		signMethod:    jwt.SigningMethodHS512,
		tokenDuration: time.Second * 60,
	}
}
func (j *jwtAuthenticator) GenerateAdminAccessToken(admin *models.AdminModel) (string, error) {
	claims := &models.AdminJwtClaims{
		Username: admin.Username,
		Role:     admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(j.signMethod, claims)

	t, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (j *jwtAuthenticator) GenerateUserAccessToken(user *models.UserModel) (string, error) {
	claims := &models.UserJwtClaims{
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(j.signMethod, claims)

	t, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (j *jwtAuthenticator) GenerateRefreshToken() (string, error) {
	uuid := uuid.New()
	return uuid.String(), nil
}

func (j *jwtAuthenticator) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func (j *jwtAuthenticator) GetClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Cannot get claims")
	}

	return claims, nil
}
