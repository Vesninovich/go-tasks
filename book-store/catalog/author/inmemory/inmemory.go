package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/stored"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Repository represents in-memory repository of authors
type Repository struct {
	data []author.StoredAuthor
	lock sync.RWMutex
}

// New creates new in-memory repository of authors
func New() *Repository {
	return &Repository{
		data: make([]author.StoredAuthor, 0),
	}
}

// GetAll gets all items from in-memory repository
func (r *Repository) GetAll(ctx context.Context) ([]book.Author, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	data := make([]book.Author, len(r.data))
	for i, s := range r.data {
		data[i] = s.ToAuthor()
	}
	return data, nil
}

// Get gets item by ID from in-memory repository
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (book.Author, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, item := range r.data {
		if item.ID == id {
			if item.IsDeleted() {
				return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
			}
			return item.ToAuthor(), nil
		}
	}

	return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
}

// Create creates item in in-memory repository
func (r *Repository) Create(ctx context.Context, dto author.CreateDTO) (book.Author, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	a := book.Author{
		ID:   uuid.New(),
		Name: dto.Name,
	}
	stored := author.StoredAuthor{
		Author: a,
		Stored: stored.Stored{
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	r.data = append(r.data, stored)
	return a, nil
}

// Update updates item in in-memory repository
func (r *Repository) Update(ctx context.Context, dto author.UpdateDTO) (book.Author, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == dto.ID {
			if item.IsDeleted() {
				return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", dto.ID)}
			}
			a := book.Author{
				ID:   dto.ID,
				Name: dto.Name,
			}
			r.data[i] = author.StoredAuthor{
				Author: a,
				Stored: stored.Stored{
					CreatedAt: item.CreatedAt,
					UpdatedAt: time.Now(),
					DeletedAt: time.Time{},
				},
			}
			return a, nil
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
			if item.IsDeleted() {
				return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
			}
			a := book.Author{
				ID:   id,
				Name: item.Name,
			}
			r.data[i] = author.StoredAuthor{
				Author: a,
				Stored: stored.Stored{
					CreatedAt: item.CreatedAt,
					UpdatedAt: item.UpdatedAt,
					DeletedAt: time.Now(),
				},
			}
			return a, nil
		}
	}
	return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
}
