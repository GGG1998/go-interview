package metrics

import "sync"

// MutexAggregator collects metrics protected by a mutex.
type MutexAggregator struct {
	mu    sync.Mutex
	count int64
	sum   float64
}

func NewMutexAggregator() *MutexAggregator {
	return &MutexAggregator{}
}

func (a *MutexAggregator) Record(value float64) {
	a.mu.Lock()
	a.count++
	a.sum += value
	a.mu.Unlock()
}

func (a *MutexAggregator) Snapshot() (count int64, sum float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.count, a.sum
}
