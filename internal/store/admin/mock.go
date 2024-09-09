package admin

import "nearbyassist/internal/models"

type MockAdminStore struct{}

func NewMockAdminStore() *MockAdminStore {
	return &MockAdminStore{}
}

func (m *MockAdminStore) Create(data *models.AdminModel) error {
	return nil
}

func (m *MockAdminStore) Login(data *models.SessionModel) error {
	return nil
}

func (m *MockAdminStore) FindById(id string) (*models.AdminModel, error) {
	return nil, nil
}

func (m *MockAdminStore) FindByUsernameHash(hash string) (*models.AdminModel, error) {
	return nil, nil
}

func (m *MockAdminStore) DoesRefreshTokenExists(refreshToken string) error {
	return nil
}

func (m *MockAdminStore) IsRefreshTokenBlacklisted(refreshToken string) error {
	return nil
}
