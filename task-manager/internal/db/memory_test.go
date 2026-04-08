package db

import (
	"slices"
	"strconv"
	"testing"
)

type IdentifiableInt int

func (i IdentifiableInt) GetId() string {
	return strconv.Itoa(int(i))
}

func TestFilterBy(t *testing.T) {
	var data map[string]IdentifiableInt = map[string]IdentifiableInt{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
	}

	db := &MemoryDb[IdentifiableInt]{
		data: data,
	}
	t.Run("Filtrowanie liczb parzystych", func(t *testing.T) {
		isEven := func(n IdentifiableInt) bool { return n%2 == 0 }
		result := slices.Collect(db.FilterBy(isEven))

		expected := []IdentifiableInt{2, 4, 6}
		if !slices.Equal(result, expected) {
			t.Errorf("Oczekiwano %v, otrzymano %v", expected, result)
		}
	})
	t.Run("Brak wyników spełniających warunek", func(t *testing.T) {
		isTooBig := func(n IdentifiableInt) bool { return n > 100 }
		result := slices.Collect(db.FilterBy(isTooBig))

		if len(result) != 0 {
			t.Errorf("Oczekiwano pustego wycinka, otrzymano %v", result)
		}
	})
}
