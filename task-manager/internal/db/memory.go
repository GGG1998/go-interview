package db

import "fmt"

type MemoryDb[T Identifiable] struct {
	data map[string]T
}

func (m *MemoryDb[T]) Insert(entity *T) (int, error) {
	if entity == nil {
		return -1, fmt.Errorf("Exist")
	}

	var e Identifiable = *entity
	m.data[e.GetId()] = *entity

	return len(m.data), nil
}

func (m *MemoryDb[T]) SelectById(id string) (*T, error) {
	if v, ok := m.data[id]; ok {
		return &v, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func NewMemoryDb[T Identifiable]() *MemoryDb[T] {
	return &MemoryDb[T]{
		data: make(map[string]T),
	}
}
