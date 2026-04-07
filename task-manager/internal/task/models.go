package task

import "time"

type Task struct {
	Id        string
	Name      string
	DueTo     string
	completed bool
}

func (t Task) GetId() string {
	return t.Id
}

/**
 * Risky solution, why?
 * Memory leak if task is older than 24H or run more than once
 * No control, how can we cancel it?
 * How it handle change dueTo?
 * */
// func (t Task) DueToAfter() <-chan time.Time {
// 	value, _ := time.Parse(time.RFC3339, t.DueTo)
// 	duration := time.Until(value)
// 	return time.After(duration)
// }

func (t Task) Duration() time.Duration {
	value, _ := time.Parse(time.RFC3339, t.DueTo)
	return time.Until(value)
}
