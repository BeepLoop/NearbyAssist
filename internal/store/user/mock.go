package user

import "nearbyassist/internal/models"

type MockUserStore struct{}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{}
}

func (m *MockUserStore) CreateUser(user *models.UserModel) (string, error) {
	return "", nil
}

func (m *MockUserStore) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	return nil, nil
}
