package category

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating Category
type CreateDTO struct {
	Name     string
	ParentID uuid.UUID
}

// UpdateDTO is DTO for updating Category
type UpdateDTO struct {
	ID       uuid.UUID
	Name     string
	ParentID uuid.UUID
}

// Repository of Categories
type Repository interface {
	GetAll(ctx context.Context) ([]book.Category, error)
	Get(ctx context.Context, id uuid.UUID) (book.Category, error)
	Create(ctx context.Context, dto CreateDTO) (book.Category, error)
	Update(ctx context.Context, dto UpdateDTO) (book.Category, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Category, error)
}
