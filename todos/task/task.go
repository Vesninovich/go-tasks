package task

import "time"

type Status uint

const (
	New Status = iota
	InProgress
	Cancelled
	Done
	Overdue
)

type Task struct {
	ID          uint64
	Name        string
	Description string
	DueDate     time.Time
	Status      Status
}
