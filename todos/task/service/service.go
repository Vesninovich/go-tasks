package taskservice

import (
	"context"
	"time"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
)

// Service handles tasks manipulation
type Service struct {
	repository task.Repository
}

// New creates new instance of Service
func New(r task.Repository) *Service {
	return &Service{r}
}

// GetAll reads all saved tasks
func (s *Service) GetAll(ctx context.Context) ([]task.Task, error) {
	return s.repository.ReadAll(ctx)
}

// CreateTask validates data, creates task if data is valid and saves it, returns error otherwise.
// All created tasks get "New" status assigned to them.
func (s *Service) CreateTask(ctx context.Context, name, desc string, dueDate int64) (task.Task, error) {
	var empty task.Task
	if name == "" {
		return empty, &common.InvalidInputError{Reason: "name is required"}
	}
	if dueDate < 0 {
		return empty, &common.InvalidInputError{Reason: "\"dueDate\" must be non-negative integer"}
	}
	due := time.Unix(dueDate, 0)
	return s.repository.Create(ctx, task.DTO{Name: name, Description: desc, DueDate: due, Status: task.New})
}
