package category

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/stored"
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

// StoredCategory is category that is stored
type StoredCategory struct {
	book.Category
	stored.Stored
}

// Repository of Categories
type Repository interface {
	GetAll(ctx context.Context) ([]book.Category, error)
	Get(ctx context.Context, id uuid.UUID) (book.Category, error)
	Create(ctx context.Context, dto CreateDTO) (book.Category, error)
	Update(ctx context.Context, dto UpdateDTO) (book.Category, error)
	Delete(ctx context.Context, id uuid.UUID) (book.Category, error)
}

// ToCategory converts stored version to actual entity
func (s StoredCategory) ToCategory() book.Category {
	return book.Category{
		ID:       s.ID,
		Name:     s.Name,
		ParentID: s.ParentID,
	}
}
