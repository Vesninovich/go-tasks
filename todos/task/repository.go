package task

import (
	"context"
	"time"
)

type DTO struct {
	Name        string
	Description string
	DueDate     time.Time
	Status      Status
}

type Repository interface {
	ReadAll(ctx context.Context) ([]Task, error)
	Create(ctx context.Context, task DTO) (Task, error)
}
