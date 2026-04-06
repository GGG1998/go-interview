package main

import (
	"net/http"

	"example.com/task-manager/internal/ratelimiter"
	"example.com/task-manager/internal/task"
)

func main() {
	mux := http.NewServeMux()

	task.NewTaskController(mux)

	mux.HandleFunc("GET /health/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("PONG"))
	})
	http.ListenAndServe(":8090", ratelimiter.RateLimiterIpMiddleware(mux))

}
