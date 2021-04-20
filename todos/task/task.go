package task

import (
	"fmt"
	"time"

	"github.com/Vesninovich/go-tasks/todos/common"
)

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

func (s Status) String() string {
	switch s {
	case New:
		return "new"
	case InProgress:
		return "in-progress"
	case Cancelled:
		return "cancelled"
	case Done:
		return "done"
	case Overdue:
		return "overdue"
	default:
		return ""
	}
}

// StatusFromText returns according enum status value
func StatusFromText(s string) (Status, error) {
	switch s {
	case "new":
		return New, nil
	case "in-progress":
		return InProgress, nil
	case "cancelled":
		return Cancelled, nil
	case "done":
		return Done, nil
	case "overdue":
		return Overdue, nil
	default:
		return New, &common.InvalidInputError{Reason: fmt.Sprintf(`status "%s" does not exist`, s)}
	}
}
