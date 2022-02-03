package author

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/stored"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating Author
type CreateDTO struct {
	Name string
}

// StoredAuthor is author that is stored
type StoredAuthor struct {
	book.Author
	stored.Stored
}

// Repository of Authors
type Repository interface {
	GetAll(ctx context.Context) ([]book.Author, error)
	Get(ctx context.Context, id uuid.UUID) (book.Author, error)
	Create(ctx context.Context, dto CreateDTO) (book.Author, error)
	Update(ctx context.Context, dto book.Author) (book.Author, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Author, error)
}

// ToAuthor converts stored version to actual entity
func (s StoredAuthor) ToAuthor() book.Author {
	return book.Author{
		ID:   s.ID,
		Name: s.Name,
	}
}
