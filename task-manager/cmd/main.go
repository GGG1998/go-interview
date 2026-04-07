package main

import (
	"net/http"
	"time"

	"example.com/task-manager/internal/db"
	"example.com/task-manager/internal/ratelimiter"
	"example.com/task-manager/internal/task"
)

func main() {
	taskDb := db.NewMemoryDb[task.Task]()
	mux := http.NewServeMux()

	task.NewTaskController(mux, taskDb)
	task.NewTaskService(taskDb)

	mux.HandleFunc("GET /health/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("PONG"))
	})

	go func() {
		http.ListenAndServe(":8090", ratelimiter.RateLimiterIpMiddleware(mux))
	}()

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			// we have to load from db task list - when fire

			// we have to load from db notification list - when fire
		}
	}()

}
