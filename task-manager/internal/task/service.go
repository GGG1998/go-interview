package task

import (
	"slices"
	"sync"
	"time"

	"example.com/task-manager/internal/db"
)

type TaskService struct {
	mu           sync.Mutex
	db           *db.MemoryDb[Task]
	activeTimers map[string]*time.Timer
}

func NewTaskService(taskDb *db.MemoryDb[Task]) *TaskService {
	return &TaskService{
		db:           taskDb,
		activeTimers: make(map[string]*time.Timer),
	}
}

func (ts *TaskService) FilterByTime(duration time.Duration) []Task {
	iters := ts.db.FilterBy(func(element Task) bool {
		return element.Duration() < duration
	})
	return slices.Collect(iters)
}

func (ts *TaskService) Schedule(id string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
}
