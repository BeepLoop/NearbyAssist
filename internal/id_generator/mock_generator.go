package id_generator

type MockGenerator struct{}

func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}

func (m *MockGenerator) Generate() (string, error) {
	return "", nil
}
