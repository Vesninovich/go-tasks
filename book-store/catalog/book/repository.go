package bookrepo

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating Book
type CreateDTO struct {
	Name       string
	Author     author.UpdateDTO
	Categories []category.UpdateDTO
}

// UpdateDTO is DTO for updating Book
type UpdateDTO struct {
	ID         uuid.UUID
	Name       string
	Author     author.UpdateDTO
	Categories []category.UpdateDTO
}

// Repository of Categories
type Repository interface {
	GetAll(ctx context.Context) ([]book.Book, error)
	Get(ctx context.Context, id uuid.UUID) (book.Book, error)
	Create(ctx context.Context, dto CreateDTO) (book.Book, error)
	Update(ctx context.Context, dto UpdateDTO) (book.Book, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Book, error)
}
