# Go Sync Benchmark: Mutex vs Channels

A benchmark suite comparing **mutex-based** and **channel-based** synchronization patterns in Go under increasing load (1,000 — 10,000 requests per second).

## Project Structure

```
cmd/                     # Executable benchmarks
  counter/               # Shared counter
  cache/                 # In-memory cache
  ratelimiter/           # Rate limiter
  workerpool/            # Worker pool
  metrics/               # Metrics aggregator
internal/                # Mutex & channel implementations
  counter/
  cache/
  ratelimiter/
  workerpool/
  metrics/
pkg/loadgen/             # Shared load generator with ramp-up
```

## How to Run

```bash
# Run a specific benchmark
go run ./cmd/counter

# Run all Go benchmarks
go test -bench=. ./...

# Run with memory profiling
go test -bench=. -benchmem ./...
```

---

## Task List

### Shared Infrastructure

- [ ] Implement load generator with configurable ramp-up (1k → 10k RPS)
- [ ] Implement latency collector (p50 / p95 / p99)
- [ ] Implement results formatter (table / CSV / JSON output)

### 1. Shared Counter

- [x] Mutex implementation
- [x] Channel implementation
- [ ] Write Go benchmarks (`_test.go`)
- [ ] Integrate with load generator for ramp-up test
- [ ] Collect and compare results

### 2. In-Memory Cache

- [x] Mutex implementation (`sync.RWMutex`)
- [x] Channel implementation
- [ ] Write Go benchmarks (`_test.go`)
- [ ] Add configurable read/write ratio (e.g. 80/20, 50/50)
- [ ] Integrate with load generator for ramp-up test
- [ ] Collect and compare results

### 3. Rate Limiter

- [x] Mutex implementation
- [x] Channel implementation
- [ ] Write Go benchmarks (`_test.go`)
- [ ] Integrate with load generator for ramp-up test
- [ ] Collect and compare results

### 4. Worker Pool

- [x] Mutex implementation (`sync.Cond` + queue)
- [x] Channel implementation (buffered channel)
- [ ] Write Go benchmarks (`_test.go`)
- [ ] Integrate with load generator for ramp-up test
- [ ] Collect and compare results

### 5. Metrics Aggregator

- [x] Mutex implementation
- [x] Channel implementation
- [ ] Write Go benchmarks (`_test.go`)
- [ ] Integrate with load generator for ramp-up test
- [ ] Collect and compare results

### Analysis & Reporting

- [ ] Generate comparison charts (latency, throughput, memory)
- [ ] Write summary of findings per benchmark
- [ ] Document when mutex wins vs when channels win
