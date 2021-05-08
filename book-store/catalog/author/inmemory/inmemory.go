package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Repository represents in-memory repository of authors
type Repository struct {
	data []book.Author
	lock sync.RWMutex
}

// New creates new in-memory repository of authors
func New() *Repository {
	return &Repository{
		data: make([]book.Author, 0),
	}
}

// GetAll gets all items from in-memory repository
func (r *Repository) GetAll(ctx context.Context) ([]book.Author, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.data[:], nil
}

// Get gets item by ID from in-memory repository
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (book.Author, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, item := range r.data {
		if item.ID == id {
			if !item.DeletedAt.IsZero() {
				return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
			}
			return item, nil
		}
	}

	return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
}

// Create creates item in in-memory repository
func (r *Repository) Create(ctx context.Context, name string) (book.Author, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	item := book.Author{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
	}
	r.data = append(r.data, item)
	return item, nil
}

// Update updates item in in-memory repository
func (r *Repository) Update(ctx context.Context, dto author.DTO) (book.Author, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if !item.DeletedAt.IsZero() {
			return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", dto.ID)}
		}
		if item.ID == dto.ID {
			r.data[i] = book.Author{
				ID:        dto.ID,
				Name:      dto.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: time.Now(),
				DeletedAt: time.Time{},
			}
			return r.data[i], nil
		}
	}
	return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", dto.ID)}
}

// Delete deletes item in in-memory repository
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Author, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == id {
			if !item.DeletedAt.IsZero() {
				return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
			}
			r.data[i] = book.Author{
				ID:        id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: time.Now(),
			}
			return r.data[i], nil
		}
	}
	return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
}
