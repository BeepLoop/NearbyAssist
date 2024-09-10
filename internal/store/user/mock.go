package user

import "nearbyassist/internal/models"

type MockUserStore struct{}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{}
}

func (m *MockUserStore) CreateUser(user *models.UserModel) (string, error) {
	return "", nil
}

func (m *MockUserStore) Login(data *models.SessionModel) error {
	return nil
}

func (m *MockUserStore) FindById(id string) (*models.UserModel, error) {
	return nil, nil
}

func (m *MockUserStore) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	return nil, nil
}

func (m *MockUserStore) DoesRefreshTokenExists(refreshToken string) error {
	return nil
}

func (m *MockUserStore) IsRefreshTokenBlacklisted(refreshToken string) error {
	return nil
}
