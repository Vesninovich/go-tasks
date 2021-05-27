package categoryservice

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Service handles categorys manipulation
type Service struct {
	repo category.Repository
}

// New creates new instance of Service
func New(r category.Repository) *Service {
	return &Service{r}
}

// GetCategory reads stored category by id
func (s *Service) GetCategory(ctx context.Context, id uuid.UUID) (book.Category, error) {
	return s.repo.Get(ctx, id)
}

// CreateCategory validates data, creates category if data is valid and saves it, returns error otherwise.
func (s *Service) CreateCategory(ctx context.Context, name string, parentID uuid.UUID) (book.Category, error) {
	var empty book.Category
	if name == "" {
		return empty, &commonerrors.InvalidInput{Reason: "name is required"}
	}
	return s.repo.Create(ctx, category.CreateDTO{Name: name, ParentID: parentID})
}
