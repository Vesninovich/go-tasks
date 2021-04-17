package task

import (
	"context"
	"time"
)

// DTO that represents Task data to be saved to storage
type DTO struct {
	Name        string
	Description string
	DueDate     time.Time
	Status      Status
}

// Repository interface represents objects that handle CRUD operations with storage
type Repository interface {
	ReadAll(ctx context.Context) ([]Task, error)
	Create(ctx context.Context, task DTO) (Task, error)
}
