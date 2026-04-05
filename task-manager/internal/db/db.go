package db

type Identifiable interface {
	GetId() string
}

type DB[T any] interface {
	Insert(*T) (int, error)
	SelectById(string) (*T, error)
}
