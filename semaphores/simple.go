package semaphores

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type Simple interface {
	Get(key string) *sync.RWMutex
}

type simple struct {
	muCache *cache.Cache
}

func NewByKey(defaultExpiration time.Duration, cleanupInterval time.Duration) Simple {
	return &simple{
		muCache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (u *simple) Get(key string) *sync.RWMutex {
	if actual, found := u.muCache.Get(key); found {
		return actual.(*sync.RWMutex)
	}

	mu := &sync.RWMutex{}
	u.muCache.Set(key, mu, cache.DefaultExpiration)
	return mu
}
