# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

A Go learning repository with hands-on exercises, tasks, and study materials covering various Go concepts.

## Structure

- **task-manager/** — exercise building an HTTP task manager API (module: `example.com/task-manager`, Go 1.25.4)
- **db-struct/** — stub for a future exercise
- **questions/** — study materials on Go concurrency (mutexes, channels) in markdown/CSV

## Commands (task-manager)

```bash
go run ./cmd          # run server on :8090
go test ./...         # all tests
go test ./internal/db/...   # single package
```

## task-manager Architecture

Layered design wired in `cmd/main.go`:

- `internal/db/` — generic `MemoryDb[T Identifiable]` with `sync.RWMutex`; `FilterBy` returns `iter.Seq[T]` (Go 1.22+)
- `internal/task/` — `Task` model, `TaskService` (scheduling with timers), `TaskController` (HTTP handlers, generic `Response[T]` JSON wrapper)
- `internal/ratelimiter/` — IP-keyed token-bucket middleware

Routes: `GET /health/`, `POST /tasks/`, `GET /tasks/{id}/`
