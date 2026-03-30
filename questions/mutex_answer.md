# Deadlocks

## Answer - What is a deadlock and why is it harder to detect than a crash?
PL:
```
Deadlock (zakleszczenie) to sytuacja, w której dwie lub więcej goroutines 
wzajemnie na siebie czekają i żadna z nich nie może kontynuować wykonania.

Jest trudniejszy do wykrycia niż crash, ponieważ aplikacja nie kończy się
z błędem — zamiast tego "cicho" przestaje przetwarzać żądania (np. timeouty API).                                                                                                                                                                        

Go wykrywa deadlock tylko wtedy, gdy WSZYSTKIE goroutines są zablokowane. - fatal error: all goroutines are asleep - deadlock!                                                                                                                                                                                 

Częściowe zakleszczenie (np. 2 z 50 goroutines) pozostaje niewykryte przez runtime. 
```

EN:
```
A deadlock occurs when two or more goroutines are waiting for each other to release resources, and none of them can proceed.

It is harder to detect than a crash because the application doesn't terminate with an error — instead, 
it silently stops processing requests (e.g. API timeouts). 

Go's runtime only detects a deadlock when ALL goroutines are blocked.
A partial deadlock (e.g. 2 out of 50 goroutines) remains undetected by the runtime.
```

## Answer - What are the four necessary conditions for a deadlock to occur?

PL:
```
- Błędna kolejność blokowania: Gdy dwie goroutines próbują zablokować te same zasoby (np. zasób A i B) w innej kolejności, co prowadzi do wzajemnego blokowania
- Podwójne blokowanie: Próba zablokowania tego samego Mutexu dwa razy przez tę samą goroutine
- Niezwolnienie blokady: Zapomnienie o wywołaniu Unlock(), co trwale blokuje dostęp innym procesom
- Pułapka RWMutex: Sytuacja, w której oczekujący "pisarz" (writer) blokuje nowych "czytelników" (readers), co przy zagnieżdżonych wywołaniach odczytu może spowodować deadlock
```

EN:
```
- Incorrect lock ordering: When two goroutines try to lock the same resources (e.g. A and B) in a different order, leading to mutual blocking
- Double locking: Attempting to lock the same Mutex twice from the same goroutine
- Forgotten unlock: Forgetting to call Unlock(), which permanently blocks access for other goroutines
- RWMutex trap: A waiting writer blocks new readers, which with nested read calls can cause a deadlock
```

## Answer - What happens if two goroutines lock resources A and B in opposite order?

PL:
```
Dochodzi do wzajemnego zakleszczenia (deadlock). Goroutine 1 blokuje zasób A i czeka na B,
a Goroutine 2 blokuje zasób B i czeka na A. Żadna z nich nie może kontynuować, bo każda
trzyma zasób, którego potrzebuje druga.

Rozwiązanie: zawsze blokuj zasoby w tej samej kolejności (np. najpierw A, potem B).
```

EN:
```
A mutual deadlock occurs. Goroutine 1 locks resource A and waits for B, while Goroutine 2
locks resource B and waits for A. Neither can proceed because each holds the resource the
other needs.

Solution: always lock resources in the same order (e.g. A first, then B).
```

## Answer - How does a deadlock occur when the same goroutine calls Lock() twice on the same mutex?

PL:
```
sync.Mutex w Go nie jest reentrant (nie wspiera zagnieżdżonego blokowania). Gdy goroutine
wywołuje Lock() na mutexie, który już trzyma — blokuje się na zawsze, czekając na Unlock(),
który nigdy nie nastąpi, bo sama jest zablokowana.
```

EN:
```
sync.Mutex in Go is not reentrant (does not support nested locking). When a goroutine calls
Lock() on a mutex it already holds — it blocks forever, waiting for an Unlock() that will
never happen because it is blocked itself.
```

## Answer - Why can forgetting Unlock() silently break an application?

PL:
```
Jeśli goroutine nie wywoła Unlock(), mutex pozostaje zablokowany na zawsze. Każda kolejna
goroutine próbująca wywołać Lock() zostanie zablokowana. Aplikacja nie crashuje — po prostu
przestaje przetwarzać żądania korzystające z tego zasobu.

Dlatego zawsze używaj defer mu.Unlock() zaraz po Lock().
```

EN:
```
If a goroutine doesn't call Unlock(), the mutex stays locked forever. Every subsequent
goroutine trying to Lock() will block. The application doesn't crash — it just silently
stops processing requests that use that resource.

Always use defer mu.Unlock() right after Lock().
```

## Answer - Why can RWMutex lead to deadlocks with nested read operations?

PL:
```
Gdy goroutine trzyma RLock() i próbuje zdobyć kolejny RLock(), a w międzyczasie inny writer
czeka na Lock() — dochodzi do deadlocka. Writer blokuje nowych czytelników, więc zagnieżdżony
RLock() nigdy nie zostanie przyznany, a pierwszy RLock() nie zostanie zwolniony.
```

EN:
```
When a goroutine holds an RLock() and tries to acquire another RLock(), while a writer is
waiting for Lock() in between — a deadlock occurs. The writer blocks new readers, so the
nested RLock() will never be granted, and the first RLock() will never be released.
```

## Answer - What happens when a writer is waiting in an RWMutex?

PL:
```
Gdy writer czeka na Lock() w RWMutex, nowi czytelnicy (RLock()) są blokowani — nawet jeśli
aktualnie inni czytelnicy jeszcze działają. To zapobiega zagłodzeniu writera (writer
starvation), ale oznacza, że oczekujący writer wstrzymuje cały nowy ruch odczytowy.
```

EN:
```
When a writer is waiting for Lock() on an RWMutex, new readers (RLock()) are blocked — even
if other readers are still active. This prevents writer starvation, but means a waiting
writer halts all new read traffic.
```

## Answer - How can you detect a deadlock in Go (tools / symptoms)?

PL:
```
- Runtime: Go wykrywa pełny deadlock automatycznie: "fatal error: all goroutines are asleep - deadlock!"
- Częściowy deadlock: Runtime go nie wykryje — objawy to timeouty API, wisząca goroutine
- go vet / staticcheck: statyczna analiza kodu
- pprof: goroutine dump (GET /debug/pprof/goroutine) pokazuje zablokowane goroutines
- Race detector: go run -race wykrywa wyścigi, choć nie bezpośrednio deadlocki
```

EN:
```
- Runtime: Go detects full deadlocks automatically: "fatal error: all goroutines are asleep - deadlock!"
- Partial deadlock: Runtime won't detect it — symptoms are API timeouts, hanging goroutines
- go vet / staticcheck: static code analysis
- pprof: goroutine dump (GET /debug/pprof/goroutine) shows blocked goroutines
- Race detector: go run -race detects races, though not deadlocks directly
```

## Answer - How would you design code to avoid deadlocks?

PL:
```
- Stała kolejność blokowania: zawsze blokuj zasoby w tej samej kolejności
- defer Unlock(): zwalniaj blokadę natychmiast po Lock()
- Minimalizuj sekcję krytyczną: trzymaj lock jak najkrócej
- Unikaj zagnieżdżonych locków: nie wywołuj Lock/RLock trzymając inny lock
- Preferuj kanały: do koordynacji złożonych potoków używaj channels zamiast mutexów
- Context z timeoutem: używaj context.WithTimeout dla operacji mogących się zablokować
```

EN:
```
- Consistent lock ordering: always lock resources in the same order
- defer Unlock(): release the lock immediately after Lock()
- Minimize critical section: hold the lock for as short as possible
- Avoid nested locks: don't call Lock/RLock while holding another lock
- Prefer channels: use channels instead of mutexes for complex pipeline coordination
- Context with timeout: use context.WithTimeout for operations that might block
```

# Performance & Contention

## Answer - What happens to sync.Mutex performance under high contention (e.g. 3400+ goroutines)?

PL:
```
Wydajność sync.Mutex drastycznie spada, gdy liczba rywalizujących goroutines przekracza
około 3400. Wynika to z narzutu mechanizmu runtime.gopark i kosztownego przełączania
kontekstu przez scheduler. 

Mutexy mogą też gorzej skalować się przy dużej liczbie rdzeni 
procesora w porównaniu do rozwiązań lock-free.
```

EN:
```
sync.Mutex performance drops drastically when the number of competing goroutines exceeds
around 3400. This is due to the overhead of the runtime.gopark mechanism and expensive
context switching by the scheduler. 

Mutexes also scale poorly with a large number of CPU
cores compared to lock-free solutions.
```

## Answer - When should you use sync/atomic instead of sync.Mutex?

PL:
```
Dla prostych operacji, takich jak liczniki, sync/atomic jest 2-3 razy szybsze niż Mutex.
Używaj sync/atomic, gdy potrzebujesz prostych operacji atomowych na pojedynczych zmiennych
(inkrementacja, odczyt, zapis). Mutex jest lepszy, gdy chronisz większą sekcję krytyczną
z wieloma operacjami.
```

EN:
```
For simple operations like counters, sync/atomic is 2-3x faster than Mutex. Use sync/atomic
when you need simple atomic operations on single variables (increment, load, store). Mutex
is better when you're protecting a larger critical section with multiple operations.
```

# Complexity & Maintainability

## Answer - Why can overusing Mutexes make Go code harder to maintain?

PL:
```
Mutexy wymagają ręcznego zarządzania blokadami, co jest podatne na błędy ludzkie — łatwo
o pominięcie Unlock() lub błędne użycie blokad odczytu/zapisu w RWMutex. Nadmierne
stosowanie Mutexów zamiast kanałów do koordynacji złożonych potoków zadań sprawia, że kod
staje się mniej czytelny i trudniejszy do zrozumienia dla innych inżynierów.
```

EN:
```
Mutexes require manual lock management, which is prone to human error — it's easy to forget
Unlock() or misuse read/write locks in RWMutex. Overusing Mutexes instead of channels for
coordinating complex task pipelines makes the code less readable and harder for other
engineers to understand.
```

# Missing Features

## Answer - Why doesn't Go's sync.Mutex support Try-Lock or timeouts?

PL:
```
sync.Mutex jest prymitywem niskopoziomowym i celowo minimalistycznym. Brakuje mu:
- Try-Lock: Brak natywnej metody sprawdzenia, czy blokada jest dostępna bez zawieszania
  goroutine (choć można to zaimplementować samodzielnie).
- Timeoutów: Brak możliwości określenia maksymalnego czasu oczekiwania na blokadę, co
  zmusza do stosowania bardziej złożonych wzorców z użyciem context i select.
```

EN:
```
sync.Mutex is a low-level, intentionally minimalist primitive. It lacks:
- Try-Lock: No native method to check if a lock is available without blocking the goroutine
  (though it can be implemented manually).
- Timeouts: No way to specify a maximum wait time for a lock, which forces developers to use
  more complex patterns with context and select.
```

# Priority Inversion

## Answer - What is priority inversion and how can it affect goroutines using Mutexes?

PL:
```
Inwersja priorytetów występuje, gdy goroutine o niskim priorytecie trzyma blokadę,
uniemożliwiając działanie krytycznym goroutines o wysokim priorytecie. Bez odpowiednich
protokołów dziedziczenia priorytetów, Mutexy mogą uniemożliwić zagwarantowanie terminów
wykonania zadań (deadlines) w systemach z ograniczeniami czasowymi.
```

EN:
```
Priority inversion occurs when a low-priority goroutine holds a lock, preventing
high-priority critical goroutines from executing. Without proper priority inheritance
protocols, Mutexes can make it impossible to guarantee task execution deadlines in
time-constrained systems.
```