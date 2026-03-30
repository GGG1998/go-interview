package ratelimiter

import (
	"sync"
	"time"
)

// MutexLimiter is a token-bucket rate limiter protected by a mutex.
type MutexLimiter struct {
	mu       sync.Mutex
	tokens   map[string]int
	limit    int
	window   time.Duration
}

func NewMutexLimiter(limit int, window time.Duration) *MutexLimiter {
	return &MutexLimiter{
		tokens: make(map[string]int),
		limit:  limit,
		window: window,
	}
}

func (l *MutexLimiter) Allow(clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens[clientID] >= l.limit {
		return false
	}
	l.tokens[clientID]++
	return true
}

func (l *MutexLimiter) Reset() {
	l.mu.Lock()
	l.tokens = make(map[string]int)
	l.mu.Unlock()
}
