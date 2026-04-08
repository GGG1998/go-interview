package db

import (
	"fmt"
	"iter"
	"sync"
)

type MemoryDb[T Identifiable] struct {
	mu   sync.RWMutex
	data map[string]T
}

func (m *MemoryDb[T]) Insert(entity *T) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entity == nil {
		return -1, fmt.Errorf("Exist")
	}

	var e Identifiable = *entity
	m.data[e.GetId()] = *entity

	return len(m.data), nil
}

func (m *MemoryDb[T]) SelectById(id string) (*T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.data[id]; ok {
		return &v, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (m *MemoryDb[T]) FilterBy(cond func(element T) bool) iter.Seq[T] {
	// m.mu.RLock()         🔴 WRONG: Block when create iterator
	// defer m.mu.RUnlock() 🔴 WRONG: Unblock immediately "return func..."
	return func(yield func(T) bool) {
		m.mu.RLock()
		defer m.mu.Unlock()
		for _, v := range m.data {
			if cond(v) && !yield(v) {
				return
			}
		}
	}
}

func NewMemoryDb[T Identifiable]() *MemoryDb[T] {
	return &MemoryDb[T]{
		data: make(map[string]T),
	}
}
