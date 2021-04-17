package task

import "time"

// Status of the Task
type Status uint

// Possible Statuses
const (
	New Status = iota
	InProgress
	Cancelled
	Done
	Overdue
)

// Task core type
type Task struct {
	ID          uint64
	Name        string
	Description string
	DueDate     time.Time
	Status      Status
}
