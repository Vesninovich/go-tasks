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
	ReadOne(ctx context.Context, id uint64) (Task, error)
	Create(ctx context.Context, task DTO) (Task, error)
	Update(ctx context.Context, id uint64, task DTO) (Task, error)
	Delete(ctx context.Context, id uint64) error
}
