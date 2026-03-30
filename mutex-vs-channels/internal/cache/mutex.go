package cache

import "sync"

// MutexCache is an in-memory key-value cache protected by RWMutex.
type MutexCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewMutexCache() *MutexCache {
	return &MutexCache{data: make(map[string]interface{})}
}

func (c *MutexCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *MutexCache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}
