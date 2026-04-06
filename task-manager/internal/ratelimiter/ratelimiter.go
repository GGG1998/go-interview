package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu         sync.Mutex
	rate       int16 // example 5
	burst      int16 // example 60
	tokens     int16
	lastUpdate time.Time
}

func NewLimiter(rate int16, burst int16) *Limiter {
	return &Limiter{
		burst:      burst,
		rate:       rate,
		tokens:     burst,
		lastUpdate: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	// JAK DAWNO ROBIŁEM REQUEST? Muszę zmierzyć róznice
	// Burst  so sub * rate
	l.mu.Lock()
	defer l.mu.Unlock()

	elapsed := int16(time.Since(l.lastUpdate).Seconds())
	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens > 1 {
		l.tokens -= 1
		return true
	}
	return false
}

// https://medium.com/@okoanton/the-internals-of-sync-map-and-its-performance-comparison-with-map-rwmutex-e000e148600c
type IpRate struct {
	mu  sync.Mutex
	ips map[string]*Limiter
}

func (ipr *IpRate) Check(ip string) bool {
	ipr.mu.Lock()
	defer ipr.mu.Unlock()

	val, ok := ipr.ips[ip]
	if !ok {
		ipr.ips[ip] = NewLimiter(5, 60)
	}

	return val.Allow()
}

func RateLimiterIpMiddleware(next http.Handler) http.Handler {
	ipRate := IpRate{
		ips: make(map[string]*Limiter),
	}

	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		if !ipRate.Check(r.RemoteAddr) {
			http.Error(response, "Limit", 403)
			return
		}
		next.ServeHTTP(response, r)
	})
}
