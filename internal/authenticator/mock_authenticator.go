package authenticator

import "github.com/golang-jwt/jwt/v5"

type MockAuthenticator struct{}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{}
}

func (m *MockAuthenticator) GenerateAdminAccessToken(options AdminOptions) (string, error) {
	return "", nil
}

func (m *MockAuthenticator) GenerateUserAccessToken(options UserOptions) (string, error) {
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
