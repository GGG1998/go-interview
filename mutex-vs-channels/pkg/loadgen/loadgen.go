package loadgen

import (
	"time"
)

// Result holds latency and error info for a single request.
type Result struct {
	Latency time.Duration
	Err     error
}

// Config defines the load generator parameters.
type Config struct {
	StartRPS   int           // starting requests per second
	EndRPS     int           // ending requests per second
	StepRPS    int           // increment per step
	StepDur    time.Duration // duration of each step
	WorkerFunc func() error  // function to call per request
}
