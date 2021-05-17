package authorservice

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Service handles authors manipulation
type Service struct {
	repo author.Repository
}

// New creates new instance of Service
func New(r author.Repository) *Service {
	return &Service{r}
}

// GetAuthor reads stored author by id
func (s *Service) GetAuthor(ctx context.Context, id uuid.UUID) (book.Author, error) {
	return s.repo.Get(ctx, id)
}

// CreateAuthor validates data, creates author if data is valid and saves it, returns error otherwise.
func (s *Service) CreateAuthor(ctx context.Context, name string) (book.Author, error) {
	var empty book.Author
	if name == "" {
		return empty, &commonerrors.InvalidInput{Reason: "name is required"}
	}
	return s.repo.Create(ctx, author.CreateDTO{Name: name})
}
