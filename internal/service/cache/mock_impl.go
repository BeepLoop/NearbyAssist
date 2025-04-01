package cache

type mockCacheImpl struct{}

func NewMockImpl() *mockCacheImpl {
	return &mockCacheImpl{}
}

func (c *mockCacheImpl) Set(key string, value interface{}) {
	return
}

func (c *mockCacheImpl) Get(key string) (interface{}, bool) {
	return nil, false
}
