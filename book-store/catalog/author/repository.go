package author

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// DTO of Author
type DTO struct {
	ID   uuid.UUID
	Name string
}

// Repository of Authors
type Repository interface {
	GetAll(ctx context.Context) ([]book.Author, error)
	Get(ctx context.Context, id uuid.UUID) (book.Author, error)
	Create(ctx context.Context, name string) (book.Author, error)
	Update(ctx context.Context, dto DTO) (book.Author, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Author, error)
}
