package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	instance *goCacheImpl
)

type goCacheImpl struct {
	store *cache.Cache
}

func NewGoCache() *goCacheImpl {
	if instance != nil {
		return instance
	}

	instance = &goCacheImpl{
		store: cache.New(time.Minute*5, time.Minute*10),
	}

	return instance
}

func (c *goCacheImpl) Set(k string, v interface{}) {
	c.store.Set(k, v, cache.DefaultExpiration)
}

func (c *goCacheImpl) Get(key string) (interface{}, bool) {
	value, found := c.store.Get(key)
	if found {
		fmt.Println("cache hit! key: ", key)
	} else {
		fmt.Println("cache miss! key: ", key)
	}

	return value, found
}
