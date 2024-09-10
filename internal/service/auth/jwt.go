package auth

import (
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ACCESS_TOKEN_ERR  = "Error generating access token"
	REFRESH_TOKEN_ERR = "Error generating refresh token"
)

type AdminJWTClaims struct {
	Id       string
	Username string
	Role     models.AdminRole
}

type UserJWTClaims struct {
	Id    string
	Name  string
	Email string
}

type Authenticator interface {
	GenerateAdminAccessToken(options AdminJWTClaims) (string, error)
	GenerateUserAccessToken(options UserJWTClaims) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) error
	GetClaims(tokenString string) (jwt.MapClaims, error)
}

type MockAuthenticator struct{}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{}
}

func (m *MockAuthenticator) GenerateAdminAccessToken(options AdminJWTClaims) (string, error) {
	return "", nil
}

func (m *MockAuthenticator) GenerateUserAccessToken(options UserJWTClaims) (string, error) {
	return "", nil
}

func (m *MockAuthenticator) GenerateRefreshToken() (string, error) {
	return "", nil
}

func (m *MockAuthenticator) ValidateToken(tokenString string) error {
	return nil
}

func (m *MockAuthenticator) GetClaims(tokenString string) (jwt.MapClaims, error) {
	return nil, nil
}

type JWTAuthenticator struct {
	secret        string
	signMethod    jwt.SigningMethod
	tokenDuration time.Duration
}

func NewJWTAuthenticator(secret string, duration int) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:        secret,
		signMethod:    jwt.SigningMethodHS512,
		tokenDuration: time.Second * time.Duration(duration),
	}
}

func (j *JWTAuthenticator) GenerateAdminAccessToken(options AdminJWTClaims) (string, error) {
	claims := &models.AdminJwtClaims{
		AdminId:  options.Id,
		Username: options.Username,
		Role:     options.Role,
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

func (j *JWTAuthenticator) GenerateUserAccessToken(options UserJWTClaims) (string, error) {
	claims := &models.UserJwtClaims{
		UserId: options.Id,
		Name:   options.Name,
		Email:  options.Email,
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

func (j *JWTAuthenticator) GenerateRefreshToken() (string, error) {
	uuid := uuid.New()
	return uuid.String(), nil
}

func (j *JWTAuthenticator) ValidateToken(tokenString string) error {
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

func (j *JWTAuthenticator) GetClaims(tokenString string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(j.secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Cannot get claims")
	}

	return claims, nil
}
