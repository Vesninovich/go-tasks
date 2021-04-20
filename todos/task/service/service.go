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

// Get reads all stored tasks
func (s *Service) Get(ctx context.Context, from, count uint) ([]task.Task, error) {
	return s.repository.Read(ctx, from, count)
}

// GetOne reads stored task by id
func (s *Service) GetOne(ctx context.Context, id uint64) (task.Task, error) {
	return s.repository.ReadOne(ctx, id)
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

// UpdateTask validates data, updates task if data is valid and task is found, returns error otherwise.
func (s *Service) UpdateTask(ctx context.Context, id uint64, name, desc string, dueDate int64, status string) (task.Task, error) {
	var empty task.Task
	if name == "" {
		return empty, &common.InvalidInputError{Reason: "name is required"}
	}
	if dueDate < 0 {
		return empty, &common.InvalidInputError{Reason: "\"dueDate\" must be non-negative integer"}
	}
	st, err := task.StatusFromText(status)
	if err != nil {
		return empty, err
	}
	due := time.Unix(dueDate, 0)
	return s.repository.Update(ctx, id, task.DTO{Name: name, Description: desc, DueDate: due, Status: st})
}

// Delete deletes stored task by id
func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
}
