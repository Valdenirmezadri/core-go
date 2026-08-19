package simplecache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
	SetWithTTL(key string, value T, d time.Duration)
	Delete(key string)
	Flush()
}

type simple[T any] struct {
	muCache *cache.Cache
}

func New[T any](defaultExpiration, cleanupInterval time.Duration) Cache[T] {
	return &simple[T]{
		muCache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (s *simple[T]) Get(key string) (T, bool) {
	actual, found := s.muCache.Get(key)
	if !found {
		var zero T
		return zero, false
	}

	return actual.(T), true
}

func (s *simple[T]) Set(key string, value T) {
	s.muCache.SetDefault(key, value)
}

// SetWithTTL an item to the cache, replacing any existing item. If the duration is 0 (DefaultExpiration), the cache's default expiration time is used. If it is -1 (NoExpiration), the item never expires.
func (s *simple[T]) SetWithTTL(key string, value T, d time.Duration) {
	s.muCache.Set(key, value, d)
}

func (s *simple[T]) Delete(key string) {
	s.muCache.Delete(key)
}

func (s *simple[T]) Flush() {
	s.muCache.Flush()
}
