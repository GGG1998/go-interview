# Deadlocks

What is a deadlock and why is it harder to detect than a crash?
What are the four necessary conditions for a deadlock to occur?
What happens if two goroutines lock resources A and B in opposite order?
How does a deadlock occur when the same goroutine calls Lock() twice on the same mutex?
Why can forgetting Unlock() silently break an application?
Why can RWMutex lead to deadlocks with nested read operations?
What happens when a writer is waiting in an RWMutex?
How can you detect a deadlock in Go (tools / symptoms)?
How would you design code to avoid deadlocks?

## Examples
Tell me how do you resolve problem?
```go
type Account struct {
    ID      int
    Balance int
    mu      sync.Mutex
}

func Transfer(from, to *Account, amount int) {
    from.mu.Lock()
    to.mu.Lock()   // Risky
    
    from.Balance -= amount
    to.Balance += amount
    
    to.mu.Unlock()
    from.mu.Unlock()
}

go Transfer(A, B, 10)
go Transfer(B, A, 10)
```

# Performance & Contention

What happens to sync.Mutex performance under high contention (e.g. 3400+ goroutines)?
When should you use sync/atomic instead of sync.Mutex?
How does RWMutex works?
Dlaczego RLock lub RWLock działają wolnej od zwykłego Locka?

# Complexity & Maintainability

Why can overusing Mutexes make Go code harder to maintain?

# Missing Features

Why doesn't Go's sync.Mutex support Try-Lock or timeouts?

# Priority Inversion

What is priority inversion and how can it affect goroutines using Mutexes?

"Sleeping" Vs "Spinning"

Metafora do mutexów:
Jesteś właścicielem stoiska piwa na festiwalu. Masz 3 krany z piwem oraz 5 pracowników.
Naturalne jest, że będziesz miał 3 kolejki z klientami. Wyobraź sobie, że pracownicy to gorutyny,
a szklanka klienta to stan, który gorutyna/nasz pracownik zmienia napełnia(pisze do niego).
Kran z piwem ma zamek(mutex) - to znaczy, że tylko jeden pracownik może napełnić kufel i zwrócić go klientowi.
Naturalne więc jest jak będzie więcej pracowników, a tylko 3 krany, to Ci pracownicy będą czekać i rywalizować, o dostęp do kranu.
A co jeśli klient też będzie gorutyną? Nikt nie będzie

Metafora do RWMutex:
Imagine that a RWMutex (Read-Write Mutex) is a door to a reading room containing a single, very important data book.
A standard Mutex is "selfish" – it only lets one person in at a time, regardless of whether they only want to read or write. 
An RWMutex is more intelligent and distinguishes between Readers and Writers.
- Multiple Readers at Once: If no one is writing, any number of people can enter and read simultaneously (RLock). This speeds up the program tremendously if you have a lot of readings.
- Only One Writer: If someone wants to write (Lock), they must be completely alone in the room. No one else can read or write at this time.
- Writer Has Priority (Queue): If someone is waiting to start writing, new readers are no longer allowed in until the writer has finished. This prevents a situation where a crowd of readers never lets the writer speak (known as writer starvation).

## Example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
	mu sync.RWMutex
	data map[string]string
}

func (db *Database) GetUser(name string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	fmt.Println("G1: Mam użytkownika, teraz sprawdzam uprawnienia...")
	
	return db.GetPermissions(name) 
}

func (db *Database) GetPermissions(name string) string {
	db.mu.RLock() 
	defer db.mu.RUnlock()
	return "Admin"
}

func main() {
	db := &Database{data: make(map[string]string)}

	go func() {
		db.GetUser("Ania")
	}()

	time.Sleep(10 * time.Millisecond)
	go func() {
		fmt.Println("Writer: I want to write, waiting for Lock()...")
		db.mu.Lock()
		fmt.Println("Writer: Write!")
		db.mu.Unlock()
	}()

	select {}
}
```