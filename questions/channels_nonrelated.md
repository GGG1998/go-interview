```go
type hchan struct {
    qcount   uint           // Ile elementów jest TERAZ w buforze
    dataqsiz uint           // Rozmiar bufora (to co podajesz w make(chan, size))
    buf      unsafe.Pointer // Wskaźnik do tablicy z danymi (bufor cykliczny)
    elemsize uint16         // Rozmiar jednego elementu (np. dla stringa to 16 bajtów)
    closed   uint32         // Flaga czy kanał jest zamknięty
    elemtype *_type        // Typ danych (string, int, itp.)
    
    sendx    uint           // Indeks w buforze, gdzie trafi następna wysłana wartość
    recvx    uint           // Indeks w buforze, skąd zostanie odebrana następna wartość

    recvq    waitq          // Kolejka zablokowanych Odbiorców (lista sudog)
    sendq    waitq          // Kolejka zablokowanych Nadawców (lista sudog)

    lock     mutex          // Blokada (kanał nie jest "magiczną" strukturą bez locków!)
}

type waitq struct {
    first *sudog
    last  *sudog
}
```

1. `make(chan string, 10)`, Go alokuje tę strukturę na stercie (heap) i zwraca Ci do niej wskaźnik
2. Każda operacja ch <- lub <- ch zaczyna się od:
    - `lock(&c.lock)`
    - Dopiero po zajęciu locka, runtime sprawdza stan bufora i kolejek. 
      To dlatego przy ekstremalnie dużej liczbie gorutyn kanały mogą stać się wąskim gardłem.
3. W polach recvq i sendq znajdują się struktury sudog.
   Gdy gorutyna musi "zasnąć", runtime nie wrzuca tam całego stosu gorutyny, 
   tylko taki "bilet parkingowy" (sudog), który zawiera:
    - Wskaźnik do samej gorutyny (żeby wiedzieć, kogo obudzić).
    - Wskaźnik do danych (żeby wiedzieć, gdzie zapisać/skąd wziąć wiadomość).

```go
// sudog (pseudo-g) represents a g in a wait list, such as for sending/receiving
// on a channel.
//
// sudog is necessary because the g ↔ synchronization object relation
// is many-to-many. A g can be on many wait lists, so there may be
// many sudogs for one g; and many gs may be waiting on the same
// synchronization object, so there may be many sudogs for one object.
//
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog

	elem maybeTraceablePtr // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	// waiters is a count of semaRoot waiting list other than head of list,
	// clamped to a uint16 to fit in unused space.
	// Only meaningful at the head of the list.
	// (If we wanted to be overly clever, we could store a high 16 bits
	// in the second entry in the list.)
	waiters uint16

	parent   *sudog             // semaRoot binary tree
	waitlink *sudog             // g.waiting list or semaRoot
	waittail *sudog             // semaRoot
	c        maybeTraceableChan // channel
}
```

Case: Dwa Odbiorniki (G1, G2) i jeden Nadawca (Main)
**Krok 1:** Gorutyna 1 wykonuje <-ch
- Gorutyna 1 chce odebrać dane, ale kanał jest pusty.
- Runtime tworzy sudog dla G1.
- Wstawia go do recvq.
- Stan waitq w recvq:
- first -> sudog (G1)
- last  -> sudog (G1)
- G1 idzie spać.

**Krok 2:** Gorutyna 2 wykonuje <-ch
- Gorutyna 2 też chce odebrać dane z tego samego kanału.
- Runtime tworzy sudog dla G2.
- Wstawia go na koniec listy recvq.
- Stan waitq w recvq:
- first -> sudog (G1)
- last  -> sudog (G2)
- (sudog G1 ma wskaźnik next na sudog G2).
- G2 idzie spać

**Krok 3:** Main wykonuje ch <- "test"
- Teraz wchodzi główny wątek z paczką danych.
- Main blokuje kanał (lock).
- Main zagląda do recvq. Widzi, że first nie jest pusty! Znajduje tam sudog Gorutyny 1.
- Direct Copy: Main kopiuje string "test" bezpośrednio do pamięci Gorutyny 1 (miejsce wskazane przez sudog.elem).
- Dequeue: Main usuwa G1 z kolejki.
- Nowy stan waitq w recvq:
- first -> sudog (G2)
- last  -> sudog (G2)
- Main budzi G1 (goready).
- Main odblokowuje kanał i leci dalej.

Dlaczego to jest bezpieczne?
Ponieważ waitq (lista sudog) jest chroniona przez hchan.lock. 
Tylko jedna gorutyna na raz może modyfikować te wskaźniki first i last. 
Dzięki temu nie ma szans, żeby dwie gorutyny "ukradły" tę samą wiadomość z kanału.