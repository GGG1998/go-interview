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

