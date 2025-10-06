package shardedbtree

import (
	"hash/fnv"
	"strings"
	"sync"
	"time"

	"github.com/google/btree"
)

const defaultDegree = 8

// item holds the cached value and expiration.
type item[T any] struct {
	value     T
	expiresAt time.Time // zero = no expiration
}

// btreeItem wraps a string key for BTree ordering.
type btreeItem struct {
	key string
}

func (a btreeItem) Less(b btree.Item) bool {
	return a.key < b.(btreeItem).key
}

// shard holds part of the cache.
type shard[T any] struct {
	mu    sync.RWMutex
	items map[string]item[T]
	index *btree.BTree // keys only
}

// ShardedBTreeCache is a thread-safe, sharded in-memory cache with TTL and prefix deletion.
type ShardedBTreeCache[T any] struct {
	shards      []*shard[T]
	numShards   uint32
	ttl         time.Duration
	cleanerStop chan struct{}
}

// NewShardedBTreeCache creates a new cache with the given number of shards and global TTL.
func NewShardedBTreeCache[T any](numShards int, globalTTL time.Duration) *ShardedBTreeCache[T] {
	c := &ShardedBTreeCache[T]{
		shards:      make([]*shard[T], numShards),
		numShards:   uint32(numShards),
		ttl:         globalTTL,
		cleanerStop: make(chan struct{}),
	}
	for i := range c.shards {
		c.shards[i] = &shard[T]{
			items: make(map[string]item[T]),
			index: btree.New(defaultDegree),
		}
	}
	go c.cleanupLoop()
	return c
}

func (c *ShardedBTreeCache[T]) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32() % c.numShards
}

func (c *ShardedBTreeCache[T]) getShard(key string) *shard[T] {
	return c.shards[c.hash(key)]
}

// Set stores a value using the global TTL.
func (c *ShardedBTreeCache[T]) Set(key string, value T) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL stores a value with a custom TTL.
func (c *ShardedBTreeCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	sh := c.getShard(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	sh.items[key] = item[T]{value: value, expiresAt: exp}
	sh.index.ReplaceOrInsert(btreeItem{key})
}

// Get retrieves a value by key.
func (c *ShardedBTreeCache[T]) Get(key string) (T, bool) {
	sh := c.getShard(key)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	var zero T
	itm, ok := sh.items[key]
	if !ok {
		return zero, false
	}
	if !itm.expiresAt.IsZero() && time.Now().After(itm.expiresAt) {
		return zero, false
	}
	return itm.value, true
}

// Delete removes a key from the cache.
func (c *ShardedBTreeCache[T]) Delete(key string) {
	sh := c.getShard(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	delete(sh.items, key)
	sh.index.Delete(btreeItem{key})
}

// ClearPrefix deletes all keys with a given prefix.
func (c *ShardedBTreeCache[T]) ClearPrefix(prefix string) {
	for _, sh := range c.shards {
		sh.mu.Lock()
		sh.index.AscendGreaterOrEqual(btreeItem{prefix}, func(i btree.Item) bool {
			k := i.(btreeItem).key
			if !strings.HasPrefix(k, prefix) {
				return false
			}
			delete(sh.items, k)
			sh.index.Delete(i)
			return true
		})
		sh.mu.Unlock()
	}
}

// cleanupLoop runs periodically to remove expired items.
func (c *ShardedBTreeCache[T]) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.cleanerStop:
			return
		}
	}
}

// cleanupExpired removes expired keys from all shards.
func (c *ShardedBTreeCache[T]) cleanupExpired() {
	now := time.Now()
	for _, sh := range c.shards {
		sh.mu.Lock()
		for k, itm := range sh.items {
			if !itm.expiresAt.IsZero() && now.After(itm.expiresAt) {
				delete(sh.items, k)
				sh.index.Delete(btreeItem{k})
			}
		}
		sh.mu.Unlock()
	}
}

// Stop stops the cleanup loop.
func (c *ShardedBTreeCache[T]) Stop() {
	close(c.cleanerStop)
}
