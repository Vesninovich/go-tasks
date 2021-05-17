package bookrepo

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/stored"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating Book
type CreateDTO struct {
	Name       string
	Author     book.Author
	Categories []book.Category
}

// StoredBook is book that is stored
type StoredBook struct {
	book.Book
	stored.Stored
}

// Repository of Categories
type Repository interface {
	Get(ctx context.Context, from, count uint, query book.Query) ([]book.Book, error)
	Create(ctx context.Context, dto CreateDTO) (book.Book, error)
	Update(ctx context.Context, dto book.Book) (book.Book, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Book, error)
}

// ToBook converts stored version to actual entity
func (s StoredBook) ToBook() book.Book {
	return book.Book{
		ID:         s.ID,
		Name:       s.Name,
		Author:     s.Author,
		Categories: s.Categories,
	}
}
