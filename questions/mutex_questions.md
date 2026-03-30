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