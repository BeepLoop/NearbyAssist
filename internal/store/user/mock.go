package user

type MockUserStore struct{}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{}
}
