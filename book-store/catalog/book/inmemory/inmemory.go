package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/stored"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Repository represents in-memory repository of books
type Repository struct {
	data []bookrepo.StoredBook
	lock sync.RWMutex
}

// New creates new in-memory repository of books
func New() *Repository {
	return &Repository{
		data: make([]bookrepo.StoredBook, 0),
	}
}

// Get fetches books
func (r *Repository) Get(ctx context.Context, from, count uint, query book.Query) ([]book.Book, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	res := make([]book.Book, 0, count)
	for _, item := range r.data {
		if item.IsDeleted() {
			continue
		}
		if matchesQuery(query, item) {
			if from > 0 {
				from--
				continue
			}
			res = append(res, book.Book{
				ID:         item.ID,
				Name:       item.Name,
				Author:     item.Author,
				Categories: item.Categories,
			})
			count--
		}
		if count == 0 {
			break
		}
	}

	return res, nil
}

func matchesQuery(query book.Query, item bookrepo.StoredBook) bool {
	if !query.Author.IsZero() && item.Author.ID != query.Author {
		return false
	}
	if len(query.Categories) > len(item.Categories) {
		return false
	}
	for _, id := range query.Categories {
		found := false
		for _, cat := range item.Categories {
			found = cat.ID == id
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Create creates item in in-memory repository
func (r *Repository) Create(ctx context.Context, dto bookrepo.CreateDTO) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	b := book.Book{
		ID:         uuid.New(),
		Name:       dto.Name,
		Author:     dto.Author,
		Categories: dto.Categories,
	}
	item := bookrepo.StoredBook{
		Book: b,
		Stored: stored.Stored{
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	r.data = append(r.data, item)
	return b, nil
}

// Update updates item in in-memory repository
func (r *Repository) Update(ctx context.Context, dto book.Book) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == dto.ID {
			if item.IsDeleted() {
				return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", dto.ID)}
			}
			b := book.Book{
				ID:         dto.ID,
				Name:       dto.Name,
				Author:     dto.Author,
				Categories: dto.Categories,
			}
			r.data[i] = bookrepo.StoredBook{
				Book: b,
				Stored: stored.Stored{
					CreatedAt: item.CreatedAt,
					UpdatedAt: time.Now(),
					DeletedAt: time.Time{},
				},
			}
			return b, nil
		}
	}
	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", dto.ID)}
}

// Delete deletes item in in-memory repository
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == id {
			if item.IsDeleted() {
				return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
			}
			b := book.Book{
				ID:         id,
				Name:       item.Name,
				Author:     item.Author,
				Categories: item.Categories,
			}
			r.data[i] = bookrepo.StoredBook{
				Book: b,
				Stored: stored.Stored{
					CreatedAt: item.CreatedAt,
					UpdatedAt: item.UpdatedAt,
					DeletedAt: time.Now(),
				},
			}
			return b, nil
		}
	}
	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
}
