package counter

import "sync"

// MutexCounter is a shared counter protected by a mutex.
type MutexCounter struct {
	mu    sync.Mutex
	value int64
}

func NewMutexCounter() *MutexCounter {
	return &MutexCounter{}
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *MutexCounter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
