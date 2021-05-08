package author

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating Author
type CreateDTO struct {
	Name string
}

// UpdateDTO is DTO for updating Author
type UpdateDTO struct {
	ID   uuid.UUID
	Name string
}

// Repository of Authors
type Repository interface {
	GetAll(ctx context.Context) ([]book.Author, error)
	Get(ctx context.Context, id uuid.UUID) (book.Author, error)
	Create(ctx context.Context, dto CreateDTO) (book.Author, error)
	Update(ctx context.Context, dto UpdateDTO) (book.Author, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Author, error)
}
