package task

import (
	"time"

	"example.com/task-manager/internal/db"
)

type TaskService struct {
	db           *db.MemoryDb[Task]
	activeTimers map[string]*time.Timer
}

func NewTaskService(taskDb *db.MemoryDb[Task]) *TaskService {
	return &TaskService{
		db:           taskDb,
		activeTimers: make(map[string]*time.Timer),
	}
}

func (ts *TaskService) Schedule(id string) {

}
