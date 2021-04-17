package task_service

import (
	"context"
	"time"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
)

type Service struct {
	repository task.Repository
}

func New(r task.Repository) *Service {
	return &Service{r}
}

func (s *Service) GetAll(ctx context.Context) ([]task.Task, error) {
	return s.repository.ReadAll(ctx)
}

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
